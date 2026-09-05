package server

import (
	"durak/internal/game"
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"math/rand/v2"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type PlayerList map[string]*Player // UserID -> Player
type Player struct {
	Name        string `json:"name"`
	UserID      string `json:"user_id"`
	Position    int    `json:"position"`
	IsConnected bool   `json:"is_connected"`

	session        Session
	JoinedAt       time.Time `json:"joined_at"`
	DisconnectedAt time.Time `json:"disconnected_at,omitempty"`
}

type LobbyList map[string]*Lobby // LobbyCode -> Lobby
type Lobby struct {
	LobbyCode string
	sync.RWMutex

	MaxPlayers int
	IsPrivate  bool
	IsPlaying  bool
	positions  []bool
	CreatedAt  time.Time

	HostID   string
	sessions map[string]string // SessionToken -> UserID
	clients  ClientList
	players  PlayerList
	game     *game.Game

	manager *Manager
}

type LobbySnapshot struct {
	LobbyCode  string    `json:"lobby_code"`
	HostID     string    `json:"host_id"`
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

type Session struct {
	LobbyCode string `json:"lobby_code"`
	Token     string `json:"token"`
}

func (m *Manager) NewLobby(userID string) *Lobby {
	l := &Lobby{
		LobbyCode:  m.GenerateLobbyCode(5),
		MaxPlayers: 4,
		IsPrivate:  false,
		IsPlaying:  false,
		HostID:     userID,
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

		host, ok := l.players[l.HostID]
		hostName := "N/A"

		if ok {
			hostName = host.Name
		}

		lobbyPreview := LobbyPreview{
			LobbyCode:   l.LobbyCode,
			HostName:    hostName,
			PlayerCount: pc,
			MaxPlayers:  l.MaxPlayers,
			CreatedAt:   l.CreatedAt,

			IsOpen:    isOpen,
			IsPlaying: l.IsPlaying,
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

func (l *Lobby) removeClient(c *Client) error {
	l.Lock()
	defer l.Unlock()

	p := l.players[c.UserID]
	if p == nil {
		return fmt.Errorf("(removeClient) Cannot find player with UserID: %s", c.UserID)
	}

	if !l.IsPlaying {
		l.freePosition(p.Position)

		c.lobby = nil

		// If host disconnects, give host to the next most recent player (that is connected.)
		if c.UserID == l.HostID {
			ps := slices.Collect(maps.Values(l.players))

			sort.Slice(ps[:], func(i, j int) bool {
				return ps[i].JoinedAt.Before(ps[j].JoinedAt)
			})

			fmt.Println(ps)

			for pl := range ps {
				player := ps[pl]
				if player.UserID != l.HostID && player.IsConnected {
					l.HostID = player.UserID
					break
				}
			}
		}

		delete(l.clients, c.UserID)
		delete(l.players, c.UserID)

		// TODO: Technically, if a player presses leave, then they leave. If they dont, then they are disconnected.
		var playerLeftMsg PlayerLeftEvent

		playerLeftMsg.UserID = c.UserID
		playerLeftMsg.HostID = l.HostID

		broadcastData, err := json.Marshal(playerLeftMsg)
		if err != nil {
			return fmt.Errorf("Failed to marshal player left message: %v", err)
		}

		playerLeft := Event{
			Payload: broadcastData,
			Type:    EventPlayerLeft,
		}

		ignored := ClientList{
			c.UserID: c,
		}
		l.broadcast(playerLeft, ignored)

		return nil
	}

	delete(l.clients, c.UserID)

	p.DisconnectedAt = time.Now()
	p.IsConnected = false

	var playerDisconnectedMsg PlayerDisconnectedEvent

	playerDisconnectedMsg.DisconnectedAt = p.DisconnectedAt
	playerDisconnectedMsg.UserID = p.UserID

	data, err := json.Marshal(playerDisconnectedMsg)
	if err != nil {
		log.Printf("Failed to marshal PlayerDisconnected message: %v", err)
		return nil
	}

	playerDisconnected := Event{
		Payload: data,
		Type:    EventPlayerDisconnected,
	}

	l.broadcast(playerDisconnected, make(ClientList))

	return nil
}

func (l *Lobby) nextAvailablePosition() int {
	for i := range l.MaxPlayers {
		if !l.positions[i] {
			return i
		}
	}

	return -1
}

func (l *Lobby) freePosition(pos int) {
	l.positions[pos] = false
}

func (l *Lobby) usePosition(pos int) int {
	l.positions[pos] = true
	return pos
}

func (l *Lobby) Snapshot() LobbySnapshot {
	snapshot := LobbySnapshot{
		LobbyCode:  l.LobbyCode,
		HostID:     l.HostID,
		Players:    l.playerSnapshots(),
		MaxPlayers: l.MaxPlayers,
		CreatedAt:  l.CreatedAt,
	}

	return snapshot
}

func (l *Lobby) SnapshotFor(c *Client) LobbySnapshot {

	snapshot := l.Snapshot()
	snapshot.Position = l.players[c.UserID].Position

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

func (l *Lobby) GenerateSessions() {
	l.sessions = make(map[string]string)

	for _, player := range l.players {
		session := Session{
			LobbyCode: l.LobbyCode,
			Token:     uuid.NewString(),
		}

		player.session = session
		l.sessions[player.UserID] = session.Token
	}

}
