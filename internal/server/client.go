package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second // how long to wait for a pong response, if 10s is up, kick

	pingInterval = (pongWait * 9) / 10 // Has to be lower than pongWait, this is the interval for the client to respond

	menuLobbiesUpdatedInterval = 10 * time.Second // TODO: Or should this be sent as a request from the client in order to allow for a timer?
)

type ClientList map[string]*Client // UserID -> Client

type Client struct {
	Name   string
	UserID string

	connection *websocket.Conn
	lobby      *Lobby
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

	pongTicker := time.NewTicker(pingInterval)
	lobbyPreviewUpdateTicker := time.NewTicker(menuLobbiesUpdatedInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok { // egress channel has closed
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed:", err)
				}
				return
			}

			if err := c.writeEvent(message); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")

		case <-pongTicker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("Error in sending ping: ", err)
				return
			}
			// fmt.Println(c.Name, "- Ping")

		case <-lobbyPreviewUpdateTicker.C:
			event, err := c.menuLobbiesUpdatedHandler()
			if err != nil {
				log.Println("Error in sending MenuLobbiesUpdated event")
				return
			}
			if event == nil {
				continue
			}

			if err := c.writeEvent(*event); err != nil {
				log.Printf("failed to send message: %v", err)
				return
			}
		}

	}
}

func (c *Client) writeEvent(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return c.connection.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) pongHandler(pongMsg string) error {
	// fmt.Println(c.Name, "- Pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) menuLobbiesUpdatedHandler() (*Event, error) {
	if c.lobby != nil {
		return nil, nil
	}

	var menuLobbiesUpdatedMsg MenuLobbiesUpdatedEvent

	menuLobbiesUpdatedMsg.Lobbies = c.manager.MenuLobbies()

	data, err := json.Marshal(menuLobbiesUpdatedMsg)
	if err != nil {
		log.Printf("Failed to marshal MenuLobbiesUpdated message: %v", err)
		return nil, err
	}

	menuLobbiesUpdated := Event{
		Payload: data,
		Type:    EventMenuLobbiesUpdated,
	}

	return &menuLobbiesUpdated, nil
}
