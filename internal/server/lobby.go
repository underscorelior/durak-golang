package server

import (
	"durak/internal/game"
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

type MemberList map[string]*Member // UserID -> Member
type Member struct {
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

	MaxMembers int
	IsHidden   bool
	InGame     bool
	positions  []bool
	CreatedAt  time.Time

	HostID   string
	sessions map[string]string // SessionToken -> UserID
	clients  ClientList
	members  MemberList
	game     *game.Game

	manager *Manager
}

type LobbySnapshot struct {
	LobbyCode  string    `json:"lobby_code"`
	HostID     string    `json:"host_id"`
	IsHidden   bool      `json:"is_hidden"`
	InGame     bool      `json:"in_game"`
	MaxMembers int       `json:"max_members"`
	CreatedAt  time.Time `json:"created_at"`

	Members []Member `json:"members"`

	Position  int            `json:"position"`
	GameState *game.Snapshot `json:"game_state,omitempty"`
}

type LobbyPreview struct {
	LobbyCode   string    `json:"lobby_code"`
	HostName    string    `json:"host_name"`
	MemberCount int       `json:"member_count"`
	MaxMembers  int       `json:"max_members"`
	CreatedAt   time.Time `json:"created_at"`

	IsOpen bool `json:"is_open"` // TODO: Need to allow user to join once started if they have a rejoin session token
	InGame bool `json:"in_game"`
}

type Session struct {
	LobbyCode    string `json:"lobby_code"`
	SessionToken string `json:"session_token"`
}

func (m *Manager) NewLobby(userID string) *Lobby {
	l := &Lobby{
		LobbyCode:  m.GenerateLobbyCode(5),
		MaxMembers: 4,
		IsHidden:   false,
		InGame:     false,
		HostID:     userID,
		CreatedAt:  time.Now(),

		clients:   make(ClientList),
		members:   make(MemberList),
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

	for _, c := range lobby.clients {
		lobby.removeClient(c)
	}

	// if _, ok := m.lobbies[lobby.LobbyCode]; ok {
	// Handle some sort of lobby deletion (DB stuff)
	delete(m.lobbies, lobby.LobbyCode)
	// }
}

func (m *Manager) MenuLobbies() []LobbyPreview {
	var lobbies []LobbyPreview

	for _, l := range m.lobbies {
		if l.IsHidden {
			continue
		}

		pc := len(l.members)

		isOpen := pc < l.MaxMembers && !l.InGame

		host, ok := l.members[l.HostID]
		hostName := "N/A"

		if ok {
			hostName = host.Name
		}

		lobbyPreview := LobbyPreview{
			LobbyCode:   l.LobbyCode,
			HostName:    hostName,
			MemberCount: pc,
			MaxMembers:  l.MaxMembers,
			CreatedAt:   l.CreatedAt,

			IsOpen: isOpen,
			InGame: l.InGame,
		}

		lobbies = append(lobbies, lobbyPreview)
	}

	sort.Slice(lobbies[:], func(i, j int) bool {
		return lobbies[i].CreatedAt.Before(lobbies[j].CreatedAt)
	})

	return lobbies
}

func (l *Lobby) addClient(c *Client, p *Member) {
	l.Lock()
	defer l.Unlock()

	l.clients[c.UserID] = c
	l.members[c.UserID] = p
}

/*
* When a client is removed:

Not during a game:
*	Free up the current position.
*	Set the lobby in the client to null
* 	If the host leaves:
* 		Either transfer to next oldest member
* 		OR close the lobby
* 	If Lobby remains open
* 	Notify new host somehow and send a message to all members that a member left

During a game:
*   Remove client from lobby
*   Send a message to all members saying the client disconnected
*/

func (l *Lobby) removeClient(c *Client) error {
	l.Lock()
	defer l.Unlock()

	p := l.members[c.UserID]
	if p == nil {
		return fmt.Errorf("(removeClient) Cannot find member with UserID: %s", c.UserID)
	}

	if !l.InGame {
		err := l.freePosition(p.Position)

		if err != nil {
			return err
		}

		c.lobby = nil

		delete(l.clients, c.UserID)
		delete(l.members, c.UserID)

		// If host disconnects while game is not active, give host to the next most recent member (that is connected.)
		if c.UserID == l.HostID {
			if len(l.members) <= 1 {
				log.Printf("(removeClient) Closing lobby %v of size %v due to host leaving\n", l.LobbyCode, len(l.members))

				l.manager.removeLobby(l)

				return nil
			}

			ps := slices.Collect(maps.Values(l.members))

			sort.Slice(ps[:], func(i, j int) bool {
				return ps[i].JoinedAt.Before(ps[j].JoinedAt)
			})

			fmt.Println(ps)

			for _, pl := range ps {
				if pl.UserID != l.HostID && pl.IsConnected {
					l.HostID = pl.UserID
					break
				}
			}
		}

		// TODO: Technically, if a member presses leave, then they leave. If they dont, then they are disconnected.
		var memberLeftPayload MemberLeftEvent
		memberLeftPayload = MemberLeftEvent{
			UserID:   c.UserID,
			HostID:   l.HostID,
			Position: p.Position,
		}

		memberLeft, err := newEvent(EventMemberLeft, memberLeftPayload)
		if err != nil {
			return err
		}

		ignored := ClientList{
			c.UserID: c,
		}
		l.broadcast(memberLeft, ignored)

		return nil
	}

	delete(l.clients, c.UserID)

	p.DisconnectedAt = time.Now()
	p.IsConnected = false

	var memberDisconnectedPayload MemberDisconnectedEvent
	memberDisconnectedPayload = MemberDisconnectedEvent{
		DisconnectedAt: p.DisconnectedAt,
		UserID:         p.UserID,
	}

	memberDisconnected, err := newEvent(EventMemberDisconnected, memberDisconnectedPayload)
	if err != nil {
		return err
	}

	l.broadcast(memberDisconnected, make(ClientList))

	return nil
}

func (l *Lobby) nextAvailablePosition() int {
	for i := range l.MaxMembers {
		if !l.positions[i] {
			return i
		}
	}

	return -1
}

func (l *Lobby) freePosition(pos int) error {
	if pos < 0 || pos >= l.MaxMembers {
		return fmt.Errorf("(freePosition) Position out of bounds")
	}

	l.positions[pos] = false

	for _, member := range l.members {
		if member.Position > pos {
			l.positions[member.Position] = false
			l.positions[member.Position-1] = true
			member.Position -= 1
		}

	}

	return nil
}

func (l *Lobby) usePosition(pos int) error {
	if pos < 0 || pos >= l.MaxMembers {
		return fmt.Errorf("(usePosition) Position out of bounds")
	}

	l.positions[pos] = true
	return nil
}

func (l *Lobby) Snapshot() LobbySnapshot {
	snapshot := LobbySnapshot{
		LobbyCode:  l.LobbyCode,
		HostID:     l.HostID,
		Members:    l.memberSnapshots(),
		MaxMembers: l.MaxMembers,
		CreatedAt:  l.CreatedAt,
		IsHidden:   l.IsHidden,
		InGame:     l.InGame,
	}

	return snapshot
}

func (l *Lobby) SnapshotFor(c *Client) *LobbySnapshot {
	p := l.members[c.UserID]

	if p == nil {
		return nil
	}

	snapshot := l.Snapshot()
	snapshot.Position = p.Position

	if l.game != nil {
		gameState := l.game.SnapshotFor(c.UserID)
		if gameState == nil {
			return nil
		}
		snapshot.GameState = gameState
	}

	return &snapshot
}

func (l *Lobby) memberSnapshots() []Member {
	var lobbyMembers []Member
	for _, p := range l.members {
		lobbyMembers = append(lobbyMembers, *p)
	}

	return lobbyMembers
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

	for _, member := range l.members {
		session := Session{
			LobbyCode:    l.LobbyCode,
			SessionToken: uuid.NewString(),
		}

		member.session = session
		l.sessions[session.SessionToken] = member.UserID
	}

}
