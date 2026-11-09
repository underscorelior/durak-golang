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

type MemberID string
type LobbyCode string
type SessionToken string

type Member struct {
	userID UserID

	Name        string   `json:"name"`
	MemberID    MemberID `json:"member_id"`
	Position    int      `json:"position"`
	IsConnected bool     `json:"is_connected"`

	session        Session
	JoinedAt       time.Time `json:"joined_at"`
	DisconnectedAt time.Time `json:"disconnected_at,omitempty"`
}
type MemberList map[UserID]*Member // UserID -> Member

type LobbyList map[LobbyCode]*Lobby // LobbyCode -> Lobby
type Lobby struct {
	Code LobbyCode
	sync.RWMutex

	MaxMembers int
	IsHidden   bool
	InGame     bool
	positions  []bool
	CreatedAt  time.Time

	hostUserID UserID
	sessions   map[SessionToken]UserID // SessionToken -> UserID
	clients    ClientList
	members    MemberList
	game       *game.Game

	manager *Manager
}

type LobbySnapshot struct {
	LobbyCode  LobbyCode `json:"lobby_code"`
	HostID     MemberID  `json:"host_id"`
	IsHidden   bool      `json:"is_hidden"`
	InGame     bool      `json:"in_game"`
	MaxMembers int       `json:"max_members"`
	CreatedAt  time.Time `json:"created_at"`

	Members []Member `json:"members"`

	Position int `json:"position"`
}

type LobbyPreview struct {
	LobbyCode   LobbyCode `json:"lobby_code"`
	HostName    string    `json:"host_name"`
	MemberCount int       `json:"member_count"`
	MaxMembers  int       `json:"max_members"`
	CreatedAt   time.Time `json:"created_at"`

	IsOpen bool `json:"is_open"` // TODO: Need to allow user to join once started if they have a rejoin session token
	InGame bool `json:"in_game"`
}

type Session struct {
	LobbyCode    LobbyCode    `json:"lobby_code"`
	SessionToken SessionToken `json:"session_token"`
}

func (m *Manager) NewLobby(userID UserID) *Lobby {
	l := &Lobby{
		Code:       m.GenerateLobbyCode(5),
		MaxMembers: 4,
		IsHidden:   false,
		InGame:     false,
		hostUserID: userID,
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

	m.lobbies[lobby.Code] = lobby
}

func (m *Manager) removeLobby(lobby *Lobby) {
	m.Lock()
	defer m.Unlock()

	for _, c := range lobby.clients {
		lobby.removeClient(c)
	}

	// if _, ok := m.lobbies[lobby.Code]; ok {
	// Handle some sort of lobby deletion (DB stuff)
	delete(m.lobbies, lobby.Code)
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

		host, ok := l.members[l.hostUserID]
		hostName := "N/A"

		if ok {
			hostName = host.Name
		}

		lobbyPreview := LobbyPreview{
			LobbyCode:   l.Code,
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

	m := l.members[c.UserID]
	if m == nil {
		return fmt.Errorf("(removeClient) Cannot find member with UserID: %s", c.UserID)
	}

	if !l.InGame {
		err := l.freePosition(m.Position)

		if err != nil {
			return err
		}

		c.lobby = nil

		delete(l.clients, c.UserID)
		delete(l.members, c.UserID)

		// If host disconnects while game is not active, give host to the next most recent member (that is connected.)
		if c.UserID == l.hostUserID {
			if len(l.members) <= 1 {
				log.Printf("(removeClient) Closing lobby %v of size %v due to host leaving\n", l.Code, len(l.members))

				l.manager.removeLobby(l)

				return nil
			}

			ps := slices.Collect(maps.Values(l.members))

			sort.Slice(ps[:], func(i, j int) bool {
				return ps[i].JoinedAt.Before(ps[j].JoinedAt)
			})

			fmt.Println(ps)

			for _, mem := range ps {
				if mem.userID != l.hostUserID && mem.IsConnected {
					l.hostUserID = mem.userID
					break
				}
			}
		}

		// TODO: Technically, if a member presses leave, then they leave. If they dont, then they are disconnected.
		var memberLeftPayload MemberLeftEvent
		memberLeftPayload = MemberLeftEvent{
			MemberID: m.MemberID,
			HostID:   l.hostMemberID(),
			Position: m.Position,
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

	m.DisconnectedAt = time.Now()
	m.IsConnected = false

	var memberDisconnectedPayload MemberDisconnectedEvent
	memberDisconnectedPayload = MemberDisconnectedEvent{
		MemberID:       m.MemberID,
		DisconnectedAt: m.DisconnectedAt,
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
		LobbyCode:  l.Code,
		HostID:     l.hostMemberID(),
		Members:    l.memberSnapshots(),
		MaxMembers: l.MaxMembers,
		CreatedAt:  l.CreatedAt,
		IsHidden:   l.IsHidden,
		InGame:     l.InGame,
	}

	return snapshot
}

func (l *Lobby) SnapshotFor(c *Client) *LobbySnapshot {
	m := l.members[c.UserID]

	if m == nil {
		return nil
	}

	snapshot := l.Snapshot()
	snapshot.Position = m.Position

	return &snapshot
}

func (l *Lobby) memberSnapshots() []Member {
	var lobbyMembers []Member
	for _, m := range l.members {
		lobbyMembers = append(lobbyMembers, *m)
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

func (m *Manager) GenerateLobbyCode(length int) LobbyCode {
	charset := strings.Split("ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")

	var code LobbyCode

	for ok := true; ok; _, ok = m.lobbies[code] {
		code = ""
		for range length {
			code += LobbyCode(charset[rand.IntN(len(charset))])
		}
	}

	return code
}

func (l *Lobby) GenerateSessions() {
	l.sessions = make(map[SessionToken]UserID)

	for _, member := range l.members {
		session := Session{
			LobbyCode:    l.Code,
			SessionToken: SessionToken(uuid.NewString()),
		}

		member.session = session
		l.sessions[session.SessionToken] = member.userID
	}

}

func (l *Lobby) hostMemberID() MemberID {
	return l.members[l.hostUserID].MemberID
}

func (m *Member) PlayerID() game.PlayerID {
	return game.PlayerID(m.MemberID)
}
