package server

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const pongWait = 10 * time.Second        // how long to wait for a pong response, if 10s is up, kick
const pingInterval = (pongWait * 9) / 10 // Has to be lower than pongWait, this is the interval for the client to respond

// TODO: Find a proper location for this struct
type ClientError struct {
	Code     string
	Message  string
	ErrorFor string
}

func (e ClientError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

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
		log.Println("(readMessages) Failed to send pong deadline for client:", c.UserID[:6])
		return
	}

	c.connection.SetReadLimit(512) // TODO: Need to know how large the largest message is

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("(readMessages) Error reading message:", err)
			}
			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Println("(readMessages) Error un-marshalling event:", err) // TODO: Manage all of these errors
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("(readMessages) Error routing event message:", err)
			break
		}
	}
}

func (c *Client) writeMessages() {
	pongTicker := time.NewTicker(pingInterval)
	defer func() {
		c.manager.removeClient(c)
		pongTicker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok { // egress channel has closed
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("(writeMessages) Connection closed:", err)
				}
				return
			}

			if err := c.writeEvent(message); err != nil {
				log.Println("(writeMessages) Failed to send message:", err)
			}

		case <-pongTicker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("(writeMessages) Error in sending ping:", err)
				// TODO: How else to handle this?
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
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) sendError(ce *ClientError, eventType string) error {
	if ce.ErrorFor == "" {
		ce.ErrorFor = eventType
	}

	payload, err := json.Marshal(ce)
	if err != nil {
		return fmt.Errorf("(sendError) Error when marshalling ErrorEvent: %w", err)
	}

	c.egress <- Event{Type: EventError, Payload: payload}
	return nil
}
