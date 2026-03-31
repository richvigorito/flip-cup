// internal/transport/api/games.go
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"flip-cup/internal/game"
	"github.com/gorilla/mux"
)

type deleteGamesResponse struct {
	Status       string   `json:"status"`
	DeletedCount int      `json:"deletedCount"`
	DeletedIDs   []string `json:"deletedIds"`
}

// GET: /games/{active|inactive|stale}
func fetchGames(manager *game.GameManager, staleAfter time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := mux.Vars(r)["status"]
		log.Printf("fetchGames called with status=%s", status)

		var response []game.GameSnapshot

		switch status {
		case "active":
			for _, g := range manager.GetAllGames() {
				if g.Active {
					response = append(response, g.Snapshot())
				}
			}
		case "inactive":
			for _, g := range manager.GetAllGames() {
				if !g.Active {
					response = append(response, g.Snapshot())
				}
			}
		case "stale":
			for _, g := range manager.GetStaleGames(staleAfter) {
				response = append(response, g.Snapshot())
			}
		default:
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}
}

// DELETE: /games/stale
func deleteStaleGames(manager *game.GameManager, staleAfter time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deletedIDs := manager.PruneStaleGames(staleAfter)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(deleteGamesResponse{
			Status:       "stale",
			DeletedCount: len(deletedIDs),
			DeletedIDs:   deletedIDs,
		})
	}
}

// DELETE: /games/{id}
func deleteGame(manager *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		log.Printf("deleting game with ID: %s", id)

		if err := manager.DeleteGameByID(id); err != nil {
			http.Error(w, "game not found or could not be deleted", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
