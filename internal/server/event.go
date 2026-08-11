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
	EventUserUpdated           = "user_updated" // "Handshake" for updating event ??? Is this best practices?
	EventJoinLobby             = "join_lobby"
	EventLobbyJoined           = "lobby_joined"
	EventLeaveLobby            = "leave_lobby"
	EventLobbyUpdate           = "lobby_updated" // Broadcasting a lobby update
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

// Sent by a client to the server to indicate joining.
type JoinLobbyEvent struct {
	Name string `json:"name"`
}

// Sent by server to recently joined client
type LobbyJoinedEvent struct {
	State game.ClientGameState
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
