package server

import (
	"durak/internal/game"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

func UpdateUser(event Event, c *Client) error {
	var updateUserEvent UpdateUserEvent
	if err := json.Unmarshal(event.Payload, &updateUserEvent); err != nil {
		return fmt.Errorf("Bad payload in request: %w", err)
	}

	c.Name = updateUserEvent.Name

	var userUpdatedMsg UserUpdatedEvent

	userUpdatedMsg.Name = c.Name

	data, err := json.Marshal(userUpdatedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal UserUpdated message: %w", err)
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
		return fmt.Errorf("Failed to marshal LobbyCreated message: %w", err)
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
		return fmt.Errorf("Bad payload in request: %w", err)
	}

	if c.lobby != nil && joinLobbyEvent.LobbyCode == c.lobby.LobbyCode {
		return ClientError{
			Code:    "ALREADY_IN_LOBBY",
			Message: "Cannot join a lobby you are already a part of",
		}
	}

	if c.lobby != nil {
		// TODO: Or should we handle this another way?
		c.lobby.removeClient(c)
	}

	lobby, ok := c.manager.lobbies[joinLobbyEvent.LobbyCode]

	if !ok {
		return ClientError{
			Code:    "LOBBY_NOT_FOUND",
			Message: fmt.Sprintf("Lobby (%v) does not exist", joinLobbyEvent.LobbyCode),
		}
	}

	if lobby.InGame {
		// TODO: Error on attempting to join while active game
		return nil
	}

	if len(lobby.players) >= int(lobby.MaxPlayers) {
		return ClientError{
			Code:    "LOBBY_FULL",
			Message: fmt.Sprintf("Lobby is full (Max %v)", lobby.MaxPlayers),
		}
	}

	if _, ok := lobby.clients[c.UserID]; ok {
		return ClientError{
			Code:    "LOBBY_DUPLICATE_PLAYER",
			Message: "Lobby contains a player with the same UserID as you",
		}
	}

	pos := lobby.nextAvailablePosition()
	if err := lobby.usePosition(pos); err != nil {
		return ClientError{
			Code:    "LOBBY_NO_VALID_POSITION",
			Message: "There was an issue finding a valid position",
		}
	}

	c.lobby = lobby

	p := &Player{
		UserID:      c.UserID,
		Name:        c.Name,
		Position:    pos,
		JoinedAt:    time.Now(),
		IsConnected: true,
	}

	lobby.addClient(c, p)

	var lobbyJoinedMsg LobbyJoinedEvent

	lobbyJoinedMsg.Lobby = lobby.SnapshotFor(c)

	data, err := json.Marshal(lobbyJoinedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal LobbyJoined message: %w", err)
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
		return fmt.Errorf("Failed to marshal PlayerJoined message: %w", err)
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

func LeaveLobby(event Event, c *Client) error {
	lobby := c.lobby

	if lobby == nil {
		// TODO: Probably send the client a message...
		return ClientError{
			Code:    "NOT_IN_LOBBY",
			Message: "Cannot leave lobby when not in lobby",
		}
	}

	if err := lobby.removeClient(c); err != nil {
		log.Println("(LeaveLobby)", err)
		return err
	}

	var lobbyLeftMsg LobbyLeftEvent

	lobbyLeftMsg.Lobbies = c.manager.MenuLobbies()

	data, err := json.Marshal(lobbyLeftMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal LobbyLeft message: %w", err)
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

	if l == nil {
		// TODO: Probably send the client a message...
		return errors.New("Cant start game when not in lobby.")
	}

	if c.UserID != l.HostID {
		// TODO: Return some sort of Error event
		return nil
	}

	if l.InGame {
		// TODO: cant start if already playing, need to return some sort of error
		return nil
	}

	if len(l.players) <= 1 {
		// TODO: Cant start if no more than 1 player
		return nil
	}

	gamePlayers := make(map[string]struct{ Position int })

	for k, v := range l.players { // TODO: Fix this BS
		gamePlayers[k] = struct{ Position int }{Position: v.Position}
	}

	l.game = game.InitializeGame(gamePlayers)
	l.InGame = true

	var gameStartedMsg GameStartedEvent
	gameStartedMsg.Lobby = l.Snapshot()

	l.GenerateSessions()

	for _, cl := range l.clients {
		gameStartedMsg.Lobby.Position = l.players[cl.UserID].Position
		gameStartedMsg.Lobby.GameState = l.game.StateFor(cl.UserID)
		gameStartedMsg.Session = l.players[cl.UserID].session

		data, err := json.Marshal(gameStartedMsg)
		if err != nil {
			return fmt.Errorf("Failed to marshal GameStarted message: %w", err)
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
		return ClientError{
			Code:    "LOBBY_NOT_FOUND",
			Message: "Lobby does not exist",
		}
	}

	userID, ok := lobby.sessions[rejoinLobbyEvent.SessionToken]

	if !ok {
		return ClientError{
			Code:    "INVALID_SESSION_TOKEN",
			Message: fmt.Sprintf("Your session token is not valid for lobby: %v", rejoinLobbyEvent.LobbyCode),
		}
	}

	// FIXME: Should we really give the newer client the older ID? Should Player have its own "lobby-specific ID", then clients have their own ID?
	if userID != c.UserID {
		return ClientError{
			Code:    "INCORRECT_USER_ID",
			Message: fmt.Sprintf("Your user ID does not match : %v", rejoinLobbyEvent.LobbyCode),
		}
	}

	p := lobby.players[c.UserID]

	if p == nil {
		return errors.New("Failed to locate player")
	}

	c.lobby = lobby
	p.IsConnected = true

	lobby.addClient(c, p)

	var lobbyRejoinedMsg LobbyRejoinedEvent

	lobbyRejoinedMsg.Lobby = lobby.SnapshotFor(c)

	data, err := json.Marshal(lobbyRejoinedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal LobbbyRejoined message: %w", err)
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
		return fmt.Errorf("Failed to marshal PlayerRejoined message: %w", err)
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

func MakeAction(event Event, c *Client) error {
	if c.lobby == nil {
		return errors.New("Cant make action when not in lobby")
	}

	if c.lobby.game == nil {
		return errors.New("Cant make an action when there is no active game")
	}

	var makeAction MakeActionEvent

	if err := json.Unmarshal(event.Payload, &makeAction); err != nil {
		return fmt.Errorf("Bad payload in request: %w", err)
	}

	if gmErr := c.lobby.game.MakeAction(c.UserID, makeAction.Action, makeAction.Card, makeAction.PairIndex); gmErr != nil {
		return gmErr
	}

	var actionMadeMsg ActionMadeEvent

	actionMadeMsg.UserID = c.UserID
	actionMadeMsg.Action = makeAction.Action
	actionMadeMsg.Phase = c.lobby.game.Turn.Phase
	actionMadeMsg.Card = makeAction.Card
	actionMadeMsg.PairIndex = makeAction.PairIndex

	actionMadeData, err := json.Marshal(actionMadeMsg)

	if err != nil {
		return fmt.Errorf("Failed to marshal ActionMade message: %w", err)
	}

	actionMade := Event{
		Payload: actionMadeData,
		Type:    EventActionMade,
	}

	c.lobby.broadcast(actionMade, ClientList{})

	return nil
}
