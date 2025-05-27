// internal/transport/api/games.go
package api

import (
	"encoding/json"
	"net/http"
	"log"
	"flip-cup/internal/game"
	"github.com/gorilla/mux"
)

//
// GET: /api/games/{active|inactive}
//
func fetchGames(manager *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)              // get path variables
		status := vars["status"]         // "active" or "inactive"
        var activeFilter *bool = nil


		// Fetch all games from the manager
        log.Println("🎯 fetchGames called") 
		games := manager.GetAllGames()


		var response = []game.GameSnapshot{}
        if status == "active" {
            activeFilter = boolPtr(true)
        } else if status == "inactive" {
            activeFilter = boolPtr(false)
        } else {
			// if you want, return 404 or 400 on invalid status
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		for _, g := range games {
			if activeFilter == nil || g.Active == *activeFilter {
				response = append(response, g.Snapshot())
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

//
// DELETE: /api/games/{active|inactive}
//
func deleteGames(manager *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)              // get path variables
		status := vars["status"]         // "active" or "inactive"
        var activeFilter *bool = nil


		// Fetch all games from the manager
        log.Println("🎯 fetchGames called") 
		games := manager.GetAllGames()


		var response = []game.GameSnapshot{}
        if status == "active" {
            activeFilter = boolPtr(true)
        } else if status == "inactive" {
            activeFilter = boolPtr(false)
        } else {
			// if you want, return 404 or 400 on invalid status
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		for _, g := range games {
			if activeFilter == nil || g.Active == *activeFilter {
                err := manager.DeleteGameByID(g.ID)
                if err != nil {
                    http.Error(w, "Game not found or could not be deleted", http.StatusNotFound)
                    return
                }
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}



//
// DELETE: /api/games/{id}
//
func deleteGame(manager *game.GameManager) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        log.Printf("Deleting game with ID: %s\n", id)

        err := manager.DeleteGameByID(id)
        if err != nil {
            http.Error(w, "Game not found or could not be deleted", http.StatusNotFound)
            return
        }

        w.WriteHeader(http.StatusNoContent) // 204 No Content on success
    }
}

func boolPtr(b bool) *bool {
    return &b
}
