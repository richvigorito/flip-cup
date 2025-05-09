package main

import (
    "fmt"
    "log"
    "net/http"
    "flip-cup/game"
    "flip-cup/ws"
)

func main() {
    // Initialize game, players, and WebSocket server
    fmt.Println("Starting game...")

    // Example of creating a game
    g := game.NewGame()

    // Load questions from YAML file
    err := g.LoadQuestions("questions.yaml")
    if err != nil {
        log.Fatalf("Failed to load questions: %v", err)
    }

    // Set up WebSocket handler
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        ws.HandleWebSocket(g, w, r) // Use the correct handler function
    })

    // Start the HTTP server
    fmt.Println("Listening on :8080")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
 
