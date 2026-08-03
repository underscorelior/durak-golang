package server

import (
	"errors"
	"log"
	"net/http"
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

type Lobby struct {
	clients ClientList
	sync.RWMutex

	CurrentID uint8 // Current ID # - used for incrememnting for user IDs

	handlers map[string]EventHandler
}

func NewLobby() *Lobby {
	l := &Lobby{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}

	l.setupEventHandlers()

	return l
}

func (l *Lobby) setupEventHandlers() {
	// l.handlers[EventJoinLobby] = ClientJoinLobby (?)
}

func (l *Lobby) routeEvent(event Event, c *Client) error {
	if handler, ok := l.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Unknown Event")
	}
}

func (l *Lobby) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	// name := r.URL.Query().Get("username")
	name := "temp"

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, l, name)

	l.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func (l *Lobby) addClient(client *Client) {
	l.Lock()
	defer l.Unlock()

	l.clients[client] = true
}

func (l *Lobby) removeClient(client *Client) {
	l.Lock()
	defer l.Unlock()

	if _, ok := l.clients[client]; ok {
		client.connection.Close()
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
