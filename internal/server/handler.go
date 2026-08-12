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
		var joinLobbyFailedMsg JoinLobbyFailedEvent

		joinLobbyFailedMsg.Code = "lobby_not_found"
		joinLobbyFailedMsg.Message = "Lobby does not exist"
		joinLobbyFailedMsg.LobbyID = joinLobbyEvent.LobbyID

		data, err := json.Marshal(joinLobbyFailedMsg)
		if err != nil {
			return fmt.Errorf("Failed to marshal JoinLobbyFailed message: %v", err)
		}

		joinLobbyFailed := Event{
			Payload: data,
			Type:    EventJoinLobbyFailed,
		}

		c.egress <- joinLobbyFailed

		return nil // I dont think this should be nil
	}

	c.lobby = lobby
	lobby.addClient(c)

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

	playerJoinedMsg.Player = LobbyPlayer{Name: c.Name, UserID: c.UserID}

	broadcastData, err := json.Marshal(playerJoinedMsg)
	if err != nil {
		return fmt.Errorf("Failed to marshal lobby joined message: %v", err)
	}

	playerJoined := Event{
		Payload: broadcastData,
		Type:    EventPlayerJoined,
	}

	ignored := ClientList{
		c: {},
	}
	lobby.broadcast(playerJoined, ignored)

	return nil
}
