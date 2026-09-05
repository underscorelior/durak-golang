package server

import (
	"encoding/json"
	"time"
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
	EventPlayerDisconnected    = "player_disconnected"
	EventRejoinLobby           = "rejoin_lobby"
	EventRejoinLobbyFailed     = "rejoin_lobby_failed"
	EventLobbyRejoined         = "lobby_rejoined"
	EventPlayerRejoined        = "player_rejoined"
	EventStartGame             = "start_game"
	EventGameStarted           = "game_started"
)

type LobbyFailedEvent struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	LobbyCode string `json:"lobby_code"`
}

type ConnectionEstablishedEvent struct {
	UserID  string         `json:"user_id"`
	Name    string         `json:"name"`
	Lobbies []LobbyPreview `json:"lobbies"`

	IsRejoining bool `json:"is_rejoining"`
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
	HostID string `json:"host_id"`
}

type PlayerDisconnectedEvent struct {
	UserID         string    `json:"user_id"`
	DisconnectedAt time.Time `json:"disconnected_at"`
}

type RejoinLobbyEvent struct {
	LobbyCode    string `json:"lobby_code"`
	SessionToken string `json:"session_token"`
}

type LobbyRejoinedEvent struct {
	Lobby LobbySnapshot `json:"lobby"`
}

type PlayerRejoinedEvent struct {
	UserID string `json:"user_id"`
}

type GameStartedEvent struct {
	Lobby        LobbySnapshot `json:"lobby"`
	SessionToken SessionToken  `json:"session_token"`
}
