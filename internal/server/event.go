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
	EventMenuLobbiesUpdated    = "menu_lobbies_updated" // TODO: later switch to a "diff" based approach, only send down stuff that changed rather than everything. also if nothing changed, send nothing at all. ALSO FIND A BETTER NAME
	EventUpdateUser            = "update_user"
	EventUserUpdated           = "user_updated"
	EventCreateLobby           = "create_lobby"
	EventLobbyCreated          = "lobby_created"
	EventJoinLobby             = "join_lobby"
	EventJoinLobbyFailed       = "join_lobby_failed"
	EventLobbyJoined           = "lobby_joined"
	EventPlayerJoined          = "player_joined"
	EventLeaveLobby            = "leave_lobby"
	EventLobbyLeft             = "lobby_left"
	EventPlayerLeft            = "player_left"
	EventStartGame             = "start_game"
	EventGameStarted           = "game_started"
)

type ConnectionEstablishedEvent struct {
	UserID  string         `json:"user_id"`
	Name    string         `json:"name"`
	Lobbies []LobbyPreview `json:"lobbies"`
}

type MenuLobbiesUpdatedEvent struct {
	Lobbies []LobbyPreview `json:"lobbies"`
}

type UpdateUserEvent struct {
	Name string `json:"name"`
}

type UserUpdatedEvent struct {
	Name string `json:"name"`
}

type LobbyCreatedEvent struct {
	LobbyCode string `json:"lobby_code"`
}

// Sent by a client to the server to indicate joining.
type JoinLobbyEvent struct {
	LobbyCode string `json:"lobby_code"`
}

type JoinLobbyFailedEvent struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	LobbyCode string `json:"lobby_code"`
}

// Sent by server to recently joined client
type LobbyJoinedEvent struct {
	Lobby LobbySnapshot `json:"lobby"`
}

type PlayerJoinedEvent struct {
	Player Player `json:"player"`
}

type LobbyLeftEvent struct {
	Lobbies []LobbyPreview `json:"lobbies"` // TODO: Maybe this isnt needed
}

type PlayerLeftEvent struct {
	UserID string `json:"user_id"`
}

type GameStartedEvent struct {
	Lobby LobbySnapshot `json:"lobby"`
}
