// cmd/flipcup/main.go
package main

import (
	"log"
	"net/http"
	"flip-cup/internal/game"
	"flip-cup/internal/transport/ws"
	"flip-cup/internal/transport/api"
	"github.com/gorilla/mux"
)

func main() {
    log.Println("✅ Flip Cup server started")
	
	manager := game.NewGameManager()

	// Create a new router
	r := mux.NewRouter()

	//  WebSocket route *  handler
	r.HandleFunc("/ws", ws.HandleWebSocketConnection(manager))

	//  HTTP routes
	api.SetupRoutes(manager, r)

	// Start the server
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", api.WithCORS(r)))
}

