package server

import (
	"durak/internal/game"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
		lobbies:  make(LobbyList),
		handlers: make(map[string]EventHandler),
	}

	m.setupEventHandlers()

	go m.broadcastMenuLobbyUpdates()

	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[EventUpdateUser] = UpdateUser
	m.handlers[EventCreateLobby] = CreateLobby
	m.handlers[EventJoinLobby] = JoinLobby
	m.handlers[EventLeaveLobby] = LeaveLobby
	m.handlers[EventStartGame] = StartGame
	m.handlers[EventRejoinLobby] = RejoinLobby
	m.handlers[EventMakeAction] = MakeAction
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	handler, ok := m.handlers[event.Type]
	if !ok {
		return c.sendError(&ClientError{Code: "UNKNOWN_EVENT", Message: fmt.Sprintf("Unknown Event: %v", event.Type)}, event.Type)
	}

	err := handler(event, c)
	if err == nil {
		return nil
	}

	var ce *ClientError
	if errors.As(err, &ce) {
		return c.sendError(ce, event.Type)
	}

	var ge *game.Error
	if errors.As(err, &ge) {
		return c.sendError(&ClientError{Code: ge.Code, Message: ge.Message}, event.Type)
	}

	log.Printf("handler %s: %v", event.Type, err)
	return c.sendError(&ClientError{Code: "INTERNAL", Message: "Something went wrong"}, event.Type)
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	name := r.URL.Query().Get("username")
	session := r.URL.Query().Get("session")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error in upgrading http connection to websocket: %v", err)
		return
	}

	var sv Session
	var client *Client
	isRejoining := false

	err = json.Unmarshal([]byte(session), &sv)
	// FIXME: might be a mistake bc err != nil could be more than just an empty string, but we will use the naive approach
	if _, ok := m.lobbies[sv.LobbyCode]; !ok || err != nil || sv.SessionToken == "" || sv.LobbyCode == "" {
		// no prior session, create a new client
		client = NewClient(conn, m, name)
	} else {
		l, lobOk := m.lobbies[sv.LobbyCode]
		// TODO: Fix whatever the hell this is
		if !lobOk {
			log.Println("!lobOk")
			client = NewClient(conn, m, name)
		} else if userId, sessOk := l.sessions[sv.SessionToken]; !sessOk {
			log.Printf("!sessOk, %v\n%v\n", sv, l.sessions)
			client = NewClient(conn, m, name)
		} else {
			log.Println("else")
			p := l.players[userId]

			client = NewClient(conn, m, p.Name)
			client.UserID = p.UserID

			isRejoining = true
		}

	}

	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()

	var connEstMsg ConnectionEstablishedEvent

	connEstMsg.Name = name
	connEstMsg.UserID = client.UserID
	connEstMsg.Lobbies = m.MenuLobbies()
	connEstMsg.IsRejoining = isRejoining

	data, err := json.Marshal(connEstMsg)
	if err != nil {
		log.Printf("Failed to marshal ConnectionEstablished message: %v", err)
		return
	}

	connectionEstablished := Event{
		Payload: data,
		Type:    EventConnectionEstablished,
	}

	client.egress <- connectionEstablished
}

func (m *Manager) addClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[c.UserID] = c
}

func (m *Manager) removeClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	if c.lobby != nil {
		c.lobby.removeClient(c)
	}

	if _, ok := m.clients[c.UserID]; ok {
		c.connection.Close()
		delete(m.clients, c.UserID)
	}
}

func (m *Manager) broadcastMenuLobbyUpdates() {
	menuLobbiesUpdatedInterval := 5 * time.Second
	lobbyPreviewUpdateTicker := time.NewTicker(menuLobbiesUpdatedInterval)

	defer func() {
		log.Println("Broadcast of MenuLobbies stopped")
		lobbyPreviewUpdateTicker.Stop()
	}()

	for range lobbyPreviewUpdateTicker.C {
		sent := 0
		var menuLobbiesUpdatedMsg MenuLobbiesUpdatedEvent

		menuLobbiesUpdatedMsg.Lobbies = m.MenuLobbies()
		menuLobbiesUpdatedMsg.NextUpdate = time.Now().Add(menuLobbiesUpdatedInterval)

		data, err := json.Marshal(menuLobbiesUpdatedMsg)
		if err != nil {
			log.Printf("Failed to marshal MenuLobbiesUpdated message: %v", err)
			continue
		}

		menuLobbiesUpdated := Event{
			Payload: data,
			Type:    EventMenuLobbiesUpdated,
		}

		for _, client := range m.clients {
			if client.lobby == nil {
				if err := client.writeEvent(menuLobbiesUpdated); err != nil {
					log.Printf("failed to send message: %v", err)
					continue
				}
				sent++
			}
		}

		log.Printf("Broadcasted %v menus to %v clients\n", len(menuLobbiesUpdatedMsg.Lobbies), sent)

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
