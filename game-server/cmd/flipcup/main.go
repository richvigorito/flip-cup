// cmd/flipcup/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"flip-cup/internal/game"
	"flip-cup/internal/transport/api"
	"flip-cup/internal/transport/ws"
	"github.com/gorilla/mux"
)

type runtimeConfig struct {
	port            string
	cleanupInterval time.Duration
	staleAfter      time.Duration
}

func main() {
	log.Println("✅ Flip Cup server started")

	config, err := loadRuntimeConfig()
	if err != nil {
		log.Fatalf("invalid runtime configuration: %v", err)
	}

	manager := game.NewGameManager()

	go func() {
		ticker := time.NewTicker(config.cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			manager.CleanupStaleGames(config.staleAfter)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/ws", ws.HandleWebSocketConnection(manager))
	api.SetupRoutes(manager, config.staleAfter, r)

	log.Printf(
		"Server running at http://localhost:%s (cleanup_interval=%s stale_after=%s)",
		config.port,
		config.cleanupInterval,
		config.staleAfter,
	)
	log.Fatal(http.ListenAndServe(":"+config.port, api.WithCORS(r)))
}

func loadRuntimeConfig() (runtimeConfig, error) {
	cleanupInterval, err := readDurationEnv("GAME_CLEANUP_INTERVAL", game.DefaultCleanupInterval)
	if err != nil {
		return runtimeConfig{}, err
	}

	staleAfter, err := readDurationEnv("GAME_STALE_AFTER", game.DefaultStaleAfter)
	if err != nil {
		return runtimeConfig{}, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return runtimeConfig{
		port:            port,
		cleanupInterval: cleanupInterval,
		staleAfter:      staleAfter,
	}, nil
}

func readDurationEnv(name string, fallback time.Duration) (time.Duration, error) {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback, nil
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", name, err)
	}

	if value <= 0 {
		return 0, fmt.Errorf("%s must be greater than zero", name)
	}

	return value, nil
}
