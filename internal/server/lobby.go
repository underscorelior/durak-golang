package server

type Lobby struct {
	clients ClientList
}

// func NewLobby() *Lobby {
// 	l := &Lobby{
// 		clients:  make(ClientList),
// 		handlers: make(map[string]EventHandler),
// 	}

// 	l.setupEventHandlers()

// 	return l
// }
