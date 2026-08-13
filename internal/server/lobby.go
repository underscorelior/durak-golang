package server

import (
	"durak/internal/game"
	"sync"

	"github.com/google/uuid"
)

type PlayerList map[string]*Player // UserID -> Player
type Player struct {
	Name     string `json:"name"`
	UserID   string `json:"userId"`
	Position int    `json:"position"`
}

type LobbyList map[string]*Lobby // LobbyID -> Lobby
type Lobby struct {
	LobbyID string
	sync.RWMutex

	MaxPlayers int
	positions  []bool

	Host    string
	clients ClientList
	players PlayerList
	game    *game.Game

	manager *Manager
}

type LobbySnapshot struct {
	LobbyID    string                  `json:"lobbyId"`
	Host       string                  `json:"host"`
	Players    []Player                `json:"players"`
	MaxPlayers int                     `json:"maxPlayers"`
	Position   int                     `json:"position"`
	GameState  *game.GameStateSnapshot `json:"gameState,omitempty"`
}

func (m *Manager) NewLobby(userID string) *Lobby {
	l := &Lobby{
		LobbyID:    uuid.NewString(),
		clients:    make(ClientList),
		players:    make(PlayerList),
		Host:       userID,
		MaxPlayers: 4,
		positions:  make([]bool, 4),
		manager:    m,
	}

	return l
}

func (m *Manager) addLobby(lobby *Lobby) {
	m.Lock()
	defer m.Unlock()

	m.lobbies[lobby.LobbyID] = lobby
}

func (m *Manager) removeLobby(lobby *Lobby) {
	m.Lock()
	defer m.Unlock()

	// if _, ok := m.lobbies[lobby.LobbyID]; ok {
	// Handle some sort of lobby deletion (DB stuff)
	delete(m.lobbies, lobby.LobbyID)
	// }
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
		LobbyID:    l.LobbyID,
		Host:       l.Host,
		Players:    l.playerSnapshots(),
		MaxPlayers: l.MaxPlayers,
		Position:   -1,
	}

	return snapshot
}

func (l *Lobby) SnapshotFor(c *Client) LobbySnapshot {
	snapshot := LobbySnapshot{
		LobbyID:    l.LobbyID,
		Host:       l.Host,
		Players:    l.playerSnapshots(),
		MaxPlayers: l.MaxPlayers,
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
