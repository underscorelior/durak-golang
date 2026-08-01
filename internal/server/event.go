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
	EventJoinLobby   = "join_lobby"
	EventLeaveLobby  = "leave_lobby"  // Leaving a lobby (client only? should i put this into a client/event file?)
	EventLobbyUpdate = "update_lobby" // Broadcasting a lobby update
)

type JoinLobbyEvent struct {
}

type LeaveLobbyEvent struct {
}

type UpdateLobbyEvent struct {
}
