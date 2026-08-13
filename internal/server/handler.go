package server

import (
	"encoding/json"
	"fmt"
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
	lobby := c.manager.NewLobby()
	c.manager.addLobby(lobby)

	var lobbyCreatedMsg LobbyCreatedEvent

	lobbyCreatedMsg.LobbyID = lobby.LobbyID

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

	if c.lobby != nil && joinLobbyEvent.LobbyID == c.lobby.LobbyID {
		return nil // How else should I handle this?
	}

	lobby, ok := c.manager.lobbies[joinLobbyEvent.LobbyID]

	if !ok {
		joinLobbyFailed, err := createLobbyFailedEvent(joinLobbyEvent.LobbyID, "lobby_not_found", "Lobby does not exist")

		if err != nil {
			return err
		}

		c.egress <- joinLobbyFailed
		return nil // I dont think this should be nil
	}

	if len(lobby.clients) >= int(lobby.MaxPlayers) {
		joinLobbyFailed, err := createLobbyFailedEvent(joinLobbyEvent.LobbyID, "lobby_full", fmt.Sprintf("Lobby is full (Max %v)", lobby.MaxPlayers))

		if err != nil {
			return err
		}

		c.egress <- joinLobbyFailed
		return nil // I dont think this should be nil
	}

	c.lobby = lobby

	pos := lobby.nextAvailablePosition()
	p := &Player{
		UserID:   c.UserID,
		Name:     c.Name,
		Position: lobby.usePosition(pos),
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
		return fmt.Errorf("Failed to marshal lobby joined message: %v", err)
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

func createLobbyFailedEvent(lobbyID, code, message string) (Event, error) {
	var joinLobbyFailedMsg JoinLobbyFailedEvent

	joinLobbyFailedMsg.Code = code
	joinLobbyFailedMsg.Message = message
	joinLobbyFailedMsg.LobbyID = lobbyID

	data, err := json.Marshal(joinLobbyFailedMsg)
	if err != nil {
		return Event{}, fmt.Errorf("Failed to marshal JoinLobbyFailed message: %v", err)
	}

	return Event{
		Payload: data,
		Type:    EventJoinLobbyFailed,
	}, nil
}
