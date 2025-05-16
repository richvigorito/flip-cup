// cmd/flipcup/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"flip-cup/internal/game"
	"flip-cup/internal/transport/ws"
	"flip-cup/internal/api"

	"github.com/gorilla/mux"
)

func main() {
	
	manager := game.NewGameManager()

	// Create a new router
	r := mux.NewRouter()

	//  HTTP routes
	api.SetupRoutes(manager, r)

	//  WebSocket route *  handler
	r.HandleFunc("/ws", ws.HandleWebSocketConnection(manager))

	// Start the server
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", api.WithCORS(r)))
}

