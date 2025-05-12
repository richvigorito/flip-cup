// main.go
package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"flip-cup/game"
	"flip-cup/ws"
	httpHandlers "flip-cup/http"
	t "flip-cup/types"

	"github.com/gorilla/mux"
)

func main() {
	manager := game.NewGameManager()

	r := mux.NewRouter()

	httpHandlers.SetupRoutes(manager, r)


	// Set up route: /ws
	r.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		// Upgrade to WebSocket connection
		conn, err := ws.Upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		// Read the first message containing gameId from the body
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading first message:", err)
			return
		}

		// Unmarshal the gameId from the incoming message
		var m t.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println("Error unmarshalling gameId:", err)
			return
		}

		log.Printf("Message %v: \n", m)

		// Get or create the game based on the gameId
		var g *game.Game
		if m.GameID != "" {
			g = manager.GetGame(m.GameID)
		}

		if g == nil {
			filename := m.Name
			if filename == "" {
    		filename = "defaultQuestion.yaml"
			}
			qf, err := t.NewQuestionFile("questions/"+filename)
			if err != nil {
				log.Fatalf("Failed to load questions: %v", err)
			}

			g = manager.NewGame(qf)
			fmt.Println("Generated new gameId:", g.ID)
		}

		// Handle the WebSocket messages for this game
		ws.HandleWebSocket(g, conn)
	})

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080",httpHandlers.WithCORS(r)))
}

