package server

import (
	"durak/internal/game"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second // how long to wait for a pong response, if 10s is up, kick

	pingInterval = (pongWait * 9) / 10 // Has to be lower than pongWait, this is the interval for the client to respond
)

type ClientList map[*Client]bool

type Client struct {
	Name    string
	UserID  string
	LobbyID *string

	gameState *game.ClientGameState

	lobby *Lobby

	connection *websocket.Conn
	manager    *Manager

	egress chan Event
}

func NewClient(conn *websocket.Conn, manager *Manager, name string) *Client {
	return &Client{
		Name:       name,
		UserID:     uuid.NewString(),
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("err")
		return
	}

	c.connection.SetReadLimit(512) // TODO: Need to know how large the largest message is

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event: %v", err) // TODO: Manage all of these errors
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handling message: ", err)
		}
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok { // egress channel has closed
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed:", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")

		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("Error in sending ping: ", err)
				return
			}
			fmt.Println(c.Name, "- Ping")
		}

	}
}

func (c *Client) pongHandler(pongMsg string) error {
	fmt.Println(c.Name, "- Pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
