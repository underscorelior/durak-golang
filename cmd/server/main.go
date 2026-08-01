package main

import (
	"log"
	"net/http"
)

func main() {

	// ctx := context.Background()

	// lobby := NewLobby(ctx)

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	// http.HandleFunc("/ws", manager.serveWS)
	// http.HandleFunc("/login", manager.loginHandler)

	// log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
