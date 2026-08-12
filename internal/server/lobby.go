package server

import (
	"sync"

	"github.com/google/uuid"
)

type LobbyList map[string]*Lobby
type Lobby struct {
	LobbyID string
	sync.RWMutex

	clients ClientList
	manager *Manager
}

func (m *Manager) NewLobby() *Lobby {
	l := &Lobby{
		clients: make(ClientList),
		LobbyID: uuid.NewString(),
		manager: m,
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

func (l *Lobby) addClient(client *Client) {
	l.Lock()
	defer l.Unlock()

	l.clients[client] = true
}

func (l *Lobby) removeClient(client *Client) {
	l.Lock()
	defer l.Unlock()

	delete(l.clients, client)
}
