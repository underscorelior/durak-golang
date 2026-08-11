package server

import (
	"encoding/json"
	"errors"
	"log"
	"maps"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex

	lobbies LobbyList

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()

	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventUpdateUser] = UpdateUser
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Unknown Event")
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	name := r.URL.Query().Get("username")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error in upgrading http connection to websocket: %v\n", err)
		return
	}

	client := NewClient(conn, m, name)

	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()

	var connEstMsg ConnectionEstablishedEvent

	connEstMsg.Name = name
	connEstMsg.Lobbies = slices.Collect(maps.Keys(m.lobbies))

	data, err := json.Marshal(connEstMsg)
	if err != nil {
		log.Printf("Failed to marshal connection established message: %v\n", err)
		return
	}

	connectionEstablished := Event{
		Payload: data,
		Type:    EventConnectionEstablished,
	}

	client.egress <- connectionEstablished
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

// CORS!
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin { // TODO: Make origin configurable from env variable
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}
