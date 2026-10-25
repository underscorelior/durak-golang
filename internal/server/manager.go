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

	log.Printf("(routeEvent) handler %s: %v\n", event.Type, err)
	return c.sendError(&ClientError{Code: "INTERNAL", Message: fmt.Sprintf("Something went wrong. %v", err.Error())}, event.Type)
}

// TODO: Clean up serveWS
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("(ServeWS) New connection")
	name := r.URL.Query().Get("username")
	session := r.URL.Query().Get("session")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("(ServeWS) Error in upgrading http connection to websocket:", err)
		return
	}

	var sv Session
	var client *Client
	isRejoining := false

	err = json.Unmarshal([]byte(session), &sv)

	// TODO: I have zero clue what is going on, rewrite this!

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
		} else if userID, sessOk := l.sessions[sv.SessionToken]; !sessOk {
			log.Printf("!sessOk, %v\n%v\n", sv, l.sessions)
			client = NewClient(conn, m, name)
		} else {
			log.Println("else")
			p := l.players[userID]

			client = NewClient(conn, m, p.Name)
			client.UserID = p.UserID

			isRejoining = true
		}

	}

	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()

	var connEstPayload ConnectionEstablishedEvent
	connEstPayload = ConnectionEstablishedEvent{
		Name:        name,
		UserID:      client.UserID,
		Lobbies:     m.MenuLobbies(),
		IsRejoining: isRejoining,
	}

	connEst, err := newEvent(EventConnectionEstablished, connEstPayload)
	if err != nil {
		err := client.sendError(&ClientError{Code: "INTERNAL", Message: fmt.Sprintf("Something went wrong. %v", err.Error())}, EventConnectionEstablished)
		fmt.Println("(ServeWS) Failed to send message:", err)
	}

	client.egress <- connEst
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
	const lobbyPrevInterval = 5 * time.Second
	lobbyPrevTicker := time.NewTicker(lobbyPrevInterval)

	defer func() {
		log.Println("Broadcast of MenuLobbies stopped.")
		lobbyPrevTicker.Stop()
	}()

	for range lobbyPrevTicker.C {
		validPlayers := 0

		for _, client := range m.clients {
			if client.lobby == nil {
				validPlayers += 1
			}
		}

		if validPlayers == 0 || len(m.clients) <= 0 {
			continue
		}

		sent := 0

		var lobbyPrevPayload MenuLobbiesUpdatedEvent
		lobbyPrevPayload.Lobbies = m.MenuLobbies()
		lobbyPrevPayload.NextUpdate = time.Now().Add(lobbyPrevInterval)

		lobbyPrev, err := newEvent(EventMenuLobbiesUpdated, lobbyPrevPayload)
		if err != nil {
			log.Println("(broadcastMenuLobbyUpdates) Failed to build event:", err)
			return
		}

		for _, client := range m.clients {
			if client.lobby == nil {
				if err := client.writeEvent(lobbyPrev); err != nil {
					log.Println("(broadcastMenuLobbyUpdates) Failed to send MenuLobbiesUpdated message:", err)
					continue
				}
				sent++
			}
		}

		log.Printf("(broadcastMenuLobbyUpdate) Broadcasted %v menus to %v clients\n", len(lobbyPrevPayload.Lobbies), sent)
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
