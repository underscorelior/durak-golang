package server

import (
	"durak/internal/game"
	"math/rand/v2"
	"sort"
	"strings"
	"sync"
	"time"
)

type PlayerList map[string]*Player // UserID -> Player
type Player struct {
	Name     string `json:"name"`
	UserID   string `json:"user_id"`
	Position int    `json:"position"`
}

type LobbyList map[string]*Lobby // LobbyCode -> Lobby
type Lobby struct {
	LobbyCode string
	sync.RWMutex

	MaxPlayers int
	IsPrivate  bool
	positions  []bool
	CreatedAt  time.Time

	Host    string
	clients ClientList
	players PlayerList
	game    *game.Game

	manager *Manager
}

type LobbySnapshot struct {
	LobbyCode  string    `json:"lobby_code"`
	Host       string    `json:"host_id"`
	IsPrivate  bool      `json:"is_private"`
	MaxPlayers int       `json:"max_players"`
	CreatedAt  time.Time `json:"created_at"`

	Players []Player `json:"players"`

	Position  int                     `json:"position"`
	GameState *game.GameStateSnapshot `json:"game_state,omitempty"`
}

type LobbyPreview struct {
	LobbyCode   string    `json:"lobby_code"`
	HostName    string    `json:"host_name"`
	PlayerCount int       `json:"player_count"`
	MaxPlayers  int       `json:"max_players"`
	CreatedAt   time.Time `json:"created_at"`

	IsOpen    bool `json:"is_open"`
	IsPlaying bool `json:"is_playing"`
}

func (m *Manager) NewLobby(userID string) *Lobby {
	l := &Lobby{
		LobbyCode:  m.GenerateLobbyCode(5),
		MaxPlayers: 4,
		IsPrivate:  false,
		Host:       userID,
		CreatedAt:  time.Now(),

		clients:   make(ClientList),
		players:   make(PlayerList),
		positions: make([]bool, 4),
		manager:   m,
	}

	return l
}

func (m *Manager) addLobby(lobby *Lobby) {
	m.Lock()
	defer m.Unlock()

	m.lobbies[lobby.LobbyCode] = lobby
}

func (m *Manager) removeLobby(lobby *Lobby) {
	m.Lock()
	defer m.Unlock()

	// if _, ok := m.lobbies[lobby.LobbyCode]; ok {
	// Handle some sort of lobby deletion (DB stuff)
	delete(m.lobbies, lobby.LobbyCode)
	// }
}

func (m *Manager) MenuLobbies() []LobbyPreview {
	var lobbies []LobbyPreview

	for lobbyCode := range m.lobbies {
		l := m.lobbies[lobbyCode]
		if l.IsPrivate {
			continue
		}

		pc := len(l.players)

		// TODO: Find a better definition of "isOpen"
		isOpen := pc < l.MaxPlayers

		lobbyPreview := LobbyPreview{
			LobbyCode: l.LobbyCode,
			// HostName:    l.players[l.Host].Name,
			HostName:    "temp",
			PlayerCount: pc,
			MaxPlayers:  l.MaxPlayers,
			CreatedAt:   l.CreatedAt,

			IsOpen:    isOpen,
			IsPlaying: l.game != nil,
		}

		lobbies = append(lobbies, lobbyPreview)
	}

	sort.Slice(lobbies[:], func(i, j int) bool {
		return lobbies[i].CreatedAt.Before(lobbies[j].CreatedAt)
	})

	return lobbies
}

func (l *Lobby) addClient(c *Client, p *Player) {
	l.Lock()
	defer l.Unlock()

	l.clients[c.UserID] = c
	l.players[c.UserID] = p
}

func (l *Lobby) removeClient(c *Client) {
	l.Lock()
	defer l.Unlock()

	delete(l.clients, c.UserID)
	delete(l.players, c.UserID)
}

func (l *Lobby) nextAvailablePosition() int {
	for i := range l.MaxPlayers {
		if !l.positions[i] {
			return i
		}
	}

	return -1
}

func (l *Lobby) usePosition(pos int) int {
	l.positions[pos] = true
	return pos
}

func (l *Lobby) Snapshot() LobbySnapshot {
	snapshot := LobbySnapshot{
		LobbyCode:  l.LobbyCode,
		Host:       l.Host,
		Players:    l.playerSnapshots(),
		MaxPlayers: l.MaxPlayers,
		Position:   -1,
	}

	return snapshot
}

func (l *Lobby) SnapshotFor(c *Client) LobbySnapshot {
	snapshot := LobbySnapshot{
		LobbyCode:  l.LobbyCode,
		Host:       l.Host,
		Players:    l.playerSnapshots(),
		MaxPlayers: l.MaxPlayers,
		CreatedAt:  l.CreatedAt,
		Position:   l.players[c.UserID].Position,
	}

	if l.game != nil {
		snapshot.GameState = l.game.StateFor(c.UserID)
	}

	return snapshot
}

func (l *Lobby) playerSnapshots() []Player {
	var lobbyPlayers []Player
	for _, p := range l.players {
		lobbyPlayers = append(lobbyPlayers, *p)
	}

	return lobbyPlayers
}

func (l *Lobby) broadcast(event Event, ignored ClientList) {
	for _, client := range l.clients {
		if _, ok := ignored[client.UserID]; ok {
			continue
		}
		client.egress <- event
	}
}

func (m *Manager) GenerateLobbyCode(length int) string {
	charset := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")

	var code string

	for ok := true; ok; _, ok = m.lobbies[code] {
		code = ""
		for range length {
			code += charset[rand.IntN(len(charset))]
		}
	}

	return code
}
