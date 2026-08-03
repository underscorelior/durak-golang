package main

import (
	"durak/internal/server"
	"log"
	"net/http"
)

func main() {

	// ctx := context.Background()

	manager := server.NewManager()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServeWS)
	// http.HandleFunc("/login", manager.loginHandler)

	// log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
