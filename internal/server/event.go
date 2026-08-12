package server

import (
	"durak/internal/game"
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
	EventLeaveLobby            = "leave_lobby"
	EventLobbyUpdate           = "lobby_updated"
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

type CreateLobbyEvent struct {
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
	LobbyID string               `json:"lobbyId"`
	State   game.ClientGameState `json:"game_state"`
}

type LeaveLobbyEvent struct {
}

// Sent out by server to all clients for stuff like when a user joins
//
// Ex: {'d_name':'player','d_op':'add','d_payload':{newplayer}}
type LobbyUpdatedEvent struct {
	Name      string `json:"delta_name"`
	Operation string `json:"delta_operation"` // Add, Remove, etc.
	Payload   any    `json:"delta_payload"`   // Replace ANY
}
