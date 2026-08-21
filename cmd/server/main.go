package main

import (
	"durak/internal/server"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	svelteURL, _ := url.Parse("http://localhost:5173")
	svelte := httputil.NewSingleHostReverseProxy(svelteURL)

	// ctx := context.Background()

	manager := server.NewManager()

	http.Handle("/", svelte)
	http.HandleFunc("/ws", manager.ServeWS)
	// http.HandleFunc("/login", manager.loginHandler)

	// log.Fatal(http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
