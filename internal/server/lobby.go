package server

import (
	"durak/internal/game"
	"sync"

	"github.com/google/uuid"
)

type LobbyList map[string]*Lobby
type Lobby struct {
	LobbyID string
	sync.RWMutex

	MaxPlayers uint8
	clients    ClientList
	game       *game.Game

	manager *Manager
}

type LobbyPlayer struct {
	Name   string `json:"name"`
	UserID string `json:"userId"`
}

type LobbySnapshot struct {
	LobbyID    string                `json:"lobbyId"`
	Players    []LobbyPlayer         `json:"players"`
	MaxPlayers uint8                 `json:"maxPlayers"`
	GameState  *game.ClientGameState `json:"gameState,omitempty"`
}

func (m *Manager) NewLobby() *Lobby {
	l := &Lobby{
		clients:    make(ClientList),
		LobbyID:    uuid.NewString(),
		MaxPlayers: 4,
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

func (l *Lobby) gameStateFor(c *Client) *game.ClientGameState {
	return nil
}

func (l *Lobby) SnapshotFor(c *Client) LobbySnapshot {
	return LobbySnapshot{
		LobbyID:    l.LobbyID,
		Players:    l.lobbyPlayers(),
		MaxPlayers: l.MaxPlayers,
		GameState:  l.gameStateFor(c),
	}
}

func (l *Lobby) lobbyPlayers() []LobbyPlayer {
	var lobbyPlayers []LobbyPlayer
	for player := range l.clients {
		lobbyPlayers = append(lobbyPlayers, LobbyPlayer{player.Name, player.UserID})
	}

	return lobbyPlayers
}
