package server

import (
	"durak/internal/game"
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventError                 = "error"
	EventConnectionEstablished = "connection_established"
	EventMenuLobbiesUpdated    = "menu_lobbies_updated" // TODO: later switch to a "diff" based approach, only send down stuff that changed rather than everything. also if nothing changed, send nothing at all. ALSO FIND A BETTER NAME
	EventUpdateUser            = "update_user"
	EventUserUpdated           = "user_updated"
	EventCreateLobby           = "create_lobby"
	EventLobbyCreated          = "lobby_created"
	EventJoinLobby             = "join_lobby"
	EventLobbyJoined           = "lobby_joined"
	EventPlayerJoined          = "player_joined"
	EventLeaveLobby            = "leave_lobby"
	EventLobbyLeft             = "lobby_left"
	EventPlayerLeft            = "player_left"
	EventPlayerDisconnected    = "player_disconnected"
	EventRejoinLobby           = "rejoin_lobby"
	EventLobbyRejoined         = "lobby_rejoined"
	EventPlayerRejoined        = "player_rejoined"
	EventStartGame             = "start_game"
	EventGameStarted           = "game_started"
	EventMakeAction            = "make_action"
	EventActionMade            = "action_made"
)

type ConnectionEstablishedEvent struct {
	UserID  string         `json:"user_id"`
	Name    string         `json:"name"`
	Lobbies []LobbyPreview `json:"lobbies"`

	IsRejoining bool `json:"is_rejoining"`
}

type MenuLobbiesUpdatedEvent struct {
	Lobbies    []LobbyPreview `json:"lobbies"`
	NextUpdate time.Time      `json:"next_update"`
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
	Lobby   LobbySnapshot `json:"lobby"`
	Session Session       `json:"session"`
}

type MakeActionEvent struct {
	Action    game.Action `json:"action"`
	Card      *game.Card  `json:"card,omitempty"`
	PairIndex *int        `json:"pair_index,omitempty"`
}

type ActionMadeEvent struct {
	UserID    string         `json:"user_id"`
	Action    game.Action    `json:"action"`
	Phase     game.TurnPhase `json:"turn_phase"`
	Card      *game.Card     `json:"card,omitempty"`
	PairIndex *int           `json:"pair_index,omitempty"`
}

type ErrorEvent struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	EventFor string `json:"event_for"`
}
