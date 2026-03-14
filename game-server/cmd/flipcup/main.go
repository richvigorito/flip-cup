// cmd/flipcup/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"flip-cup/internal/game"
	"flip-cup/internal/transport/ws"
	"flip-cup/internal/transport/api"
	"github.com/gorilla/mux"
)

func main() {
    log.Println("✅ Flip Cup server started")
	
	manager := game.NewGameManager()

	// Start background cleanup task
	go func() {
		// Run cleanup every 30 minutes, removing games inactive for > 1 hour
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			manager.CleanupStaleGames(1 * time.Hour)
		}
	}()

	// Create a new router
	r := mux.NewRouter()

	//  WebSocket route *  handler
	r.HandleFunc("/ws", ws.HandleWebSocketConnection(manager))

	//  HTTP routes
	api.SetupRoutes(manager, r)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, api.WithCORS(r)))
}

