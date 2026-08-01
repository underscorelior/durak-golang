package server

import (
	"durak/internal/client"
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
	clients client.ClientList
	sync.RWMutex
}

// CORS!
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin { // Make origin configurable from env variable
	case "https://localhost:8080":
		return true
	default:
		return false
	}
}
