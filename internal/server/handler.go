package server

import (
	"durak/internal/game"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"
	"time"
)

func UpdateUser(event Event, c *Client) error {
	var updateUserEvent UpdateUserEvent
	if err := json.Unmarshal(event.Payload, &updateUserEvent); err != nil {
		return fmt.Errorf("Bad payload in request: %v", err)
	}

	c.Name = updateUserEvent.Name

	var userUpdatedMsg UserUpdatedEvent

	userUpdatedMsg.Name = c.Name

	data, err := json.Marshal(userUpdatedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal UserUpdated message: %v", err)
	}

	userUpdated := Event{
		Payload: data,
		Type:    EventUserUpdated,
	}

	c.egress <- userUpdated

	return nil
}

func CreateLobby(event Event, c *Client) error {
	// TODO: Need to implement ratelimit system
	lobby := c.manager.NewLobby(c.UserID)
	c.manager.addLobby(lobby)

	var lobbyCreatedMsg LobbyCreatedEvent

	lobbyCreatedMsg.LobbyCode = lobby.LobbyCode

	data, err := json.Marshal(lobbyCreatedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal LobbyCreated message: %v", err)
	}

	lobbyCreated := Event{
		Payload: data,
		Type:    EventLobbyCreated,
	}

	c.egress <- lobbyCreated

	return nil
}

func JoinLobby(event Event, c *Client) error {
	var joinLobbyEvent JoinLobbyEvent
	if err := json.Unmarshal(event.Payload, &joinLobbyEvent); err != nil {
		return fmt.Errorf("Bad payload in request: %v", err)
	}

	if c.lobby != nil && joinLobbyEvent.LobbyCode == c.lobby.LobbyCode {
		return nil // How else should I handle this?
	}

	lobby, ok := c.manager.lobbies[joinLobbyEvent.LobbyCode]

	if !ok {
		joinLobbyFailed, err := createLobbyFailedEvent(joinLobbyEvent.LobbyCode, "lobby_not_found", "Lobby does not exist", EventJoinLobbyFailed)

		if err != nil {
			return err
		}

		c.egress <- joinLobbyFailed
		return nil // I dont think this should be nil
	}

	if len(lobby.players) >= int(lobby.MaxPlayers) {
		joinLobbyFailed, err := createLobbyFailedEvent(joinLobbyEvent.LobbyCode, "lobby_full", fmt.Sprintf("Lobby is full (Max %v)", lobby.MaxPlayers), EventJoinLobbyFailed)

		if err != nil {
			return err
		}

		c.egress <- joinLobbyFailed
		return nil // I dont think this should be nil
	}

	if _, ok := lobby.clients[c.UserID]; ok {
		joinLobbyFailed, err := createLobbyFailedEvent(joinLobbyEvent.LobbyCode, "lobby_duplicate_player", "Lobby contains a player with the same UserID as you", EventJoinLobbyFailed)

		if err != nil {
			return err
		}

		c.egress <- joinLobbyFailed
		return nil // I dont think this should be nil
	}

	c.lobby = lobby

	pos := lobby.nextAvailablePosition()
	p := &Player{
		UserID:      c.UserID,
		Name:        c.Name,
		Position:    lobby.usePosition(pos),
		JoinedAt:    time.Now(),
		IsConnected: true,
	}

	lobby.addClient(c, p)

	var lobbyJoinedMsg LobbyJoinedEvent

	lobbyJoinedMsg.Lobby = lobby.SnapshotFor(c)

	data, err := json.Marshal(lobbyJoinedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal lobby joined message: %v", err)
	}

	lobbyJoined := Event{
		Payload: data,
		Type:    EventLobbyJoined,
	}

	c.egress <- lobbyJoined

	var playerJoinedMsg PlayerJoinedEvent

	playerJoinedMsg.Player = *p

	broadcastData, err := json.Marshal(playerJoinedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal player joined message: %v", err)
	}

	playerJoined := Event{
		Payload: broadcastData,
		Type:    EventPlayerJoined,
	}

	ignored := ClientList{
		c.UserID: c,
	}
	lobby.broadcast(playerJoined, ignored)

	return nil
}

func createLobbyFailedEvent(lobbyCode, code, message, evtType string) (Event, error) {
	var joinLobbyFailedMsg LobbyFailedEvent

	joinLobbyFailedMsg.Code = code
	joinLobbyFailedMsg.Message = message
	joinLobbyFailedMsg.LobbyCode = lobbyCode

	data, err := json.Marshal(joinLobbyFailedMsg)
	if err != nil {
		return Event{}, fmt.Errorf("Failed to marshal JoinLobbyFailed message: %v", err)
	}

	return Event{
		Payload: data,
		Type:    evtType,
	}, nil
}

func LeaveLobby(event Event, c *Client) error {
	lobby := c.lobby

	if err := lobby.removeClient(c); err != nil {
		log.Println(err)
		return err
	}

	var lobbyLeftMsg LobbyLeftEvent

	lobbyLeftMsg.Lobbies = c.manager.MenuLobbies()

	data, err := json.Marshal(lobbyLeftMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal lobby left message: %v", err)
	}

	lobbyLeft := Event{
		Payload: data,
		Type:    EventLobbyLeft,
	}

	c.egress <- lobbyLeft

	return nil
}

func StartGame(event Event, c *Client) error {
	l := c.lobby

	if c.UserID != l.HostID {
		// TODO: Return some sort of Error event
		return nil
	}

	l.game = game.InitializeGame(slices.Collect(maps.Keys(l.players)))
	l.IsPlaying = true

	var gameStartedMsg GameStartedEvent
	gameStartedMsg.Lobby = l.Snapshot()

	l.GenerateSessions()

	for _, cl := range c.lobby.clients {
		gameStartedMsg.Lobby.Position = l.players[cl.UserID].Position
		gameStartedMsg.Lobby.GameState = l.game.StateFor(cl.UserID)
		gameStartedMsg.Session = l.players[cl.UserID].session

		data, err := json.Marshal(gameStartedMsg)
		if err != nil {
			return fmt.Errorf("Failed to marshal lobby joined message: %v", err)
		}

		event := Event{
			Payload: data,
			Type:    EventGameStarted,
		}
		cl.egress <- event
	}

	return nil
}

func RejoinLobby(event Event, c *Client) error {
	var rejoinLobbyEvent RejoinLobbyEvent

	if err := json.Unmarshal(event.Payload, &rejoinLobbyEvent); err != nil {
		return fmt.Errorf("Bad payload in request: %v", err)
	}

	if c.lobby != nil && rejoinLobbyEvent.LobbyCode == c.lobby.LobbyCode {
		return nil // How else should I handle this?
	}

	lobby, ok := c.manager.lobbies[rejoinLobbyEvent.LobbyCode]

	if !ok {
		rejoinLobbyFailed, err := createLobbyFailedEvent(rejoinLobbyEvent.LobbyCode, "lobby_not_found", "Lobby does not exist", EventRejoinLobbyFailed)

		if err != nil {
			return err
		}

		c.egress <- rejoinLobbyFailed
		return nil // I dont think this should be nil
	}

	userID, ok := lobby.sessions[rejoinLobbyEvent.SessionToken]

	if !ok {
		rejoinLobbyFailed, err := createLobbyFailedEvent(rejoinLobbyEvent.LobbyCode, "invalid_session_token", fmt.Sprintf("Your session token is not valid for lobby: %v", rejoinLobbyEvent.LobbyCode), EventRejoinLobbyFailed)

		if err != nil {
			return err
		}

		c.egress <- rejoinLobbyFailed
		return nil // I dont think this should be nil
	}

	// FIXME: Should we really give the newer client the older ID? Should Player have its own "lobby-specific ID", then clients have their own ID?
	if userID != c.UserID {
		rejoinLobbyFailed, err := createLobbyFailedEvent(rejoinLobbyEvent.LobbyCode, "incorrect_user_id", fmt.Sprintf("Your user ID does not match : %v", rejoinLobbyEvent.LobbyCode), EventRejoinLobbyFailed)

		if err != nil {
			return err
		}

		c.egress <- rejoinLobbyFailed
		return nil // I dont think this should be nil
	}

	c.lobby = lobby
	p := lobby.players[c.UserID]

	p.IsConnected = true

	lobby.addClient(c, p)

	var lobbyRejoinedMsg LobbyRejoinedEvent

	lobbyRejoinedMsg.Lobby = lobby.SnapshotFor(c)

	data, err := json.Marshal(lobbyRejoinedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal lobby rejoined message: %v", err)
	}

	lobbyRejoined := Event{
		Payload: data,
		Type:    EventLobbyRejoined,
	}

	c.egress <- lobbyRejoined

	var playerRejoinedMsg PlayerRejoinedEvent

	playerRejoinedMsg.UserID = c.UserID

	broadcastData, err := json.Marshal(playerRejoinedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal player joined message: %v", err)
	}

	playerRejoined := Event{
		Payload: broadcastData,
		Type:    EventPlayerRejoined,
	}

	ignored := ClientList{
		c.UserID: c,
	}
	lobby.broadcast(playerRejoined, ignored)

	return nil

}
