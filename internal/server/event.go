package server

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventConnectionEstablished = "connection_established"
	EventUpdateUser            = "update_user"
	EventUserUpdated           = "user_updated"
	EventCreateLobby           = "create_lobby"
	EventLobbyCreated          = "lobby_created"
	EventJoinLobby             = "join_lobby"
	EventJoinLobbyFailed       = "join_lobby_failed"
	EventLobbyJoined           = "lobby_joined"
	EventPlayerJoined          = "player_joined"
	EventLeaveLobby            = "leave_lobby"
	EventPlayerLeft            = "player_left"
)

type ConnectionEstablishedEvent struct {
	Name    string   `json:"name"`
	Lobbies []string `json:"lobbies"` // I think I need to make this to be a json string, will figure out later
}

type UpdateUserEvent struct {
	Name string `json:"name"`
}

type UserUpdatedEvent struct {
	Name string `json:"name"`
}

type LobbyCreatedEvent struct {
	LobbyID string `json:"lobbyId"`
}

// Sent by a client to the server to indicate joining.
type JoinLobbyEvent struct {
	LobbyID string `json:"lobbyId"`
}

type JoinLobbyFailedEvent struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	LobbyID string `json:"lobbyId"`
}

// Sent by server to recently joined client
type LobbyJoinedEvent struct {
	Lobby LobbySnapshot `json:"lobby"`
}

type PlayerJoinedEvent struct {
	Player Player `json:"player"`
}
