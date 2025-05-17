// internal/api/router.go
package api

import (
	"encoding/json"
	"net/http"
	"log"
	"flip-cup/internal/game"
)

func fetchGames(manager *game.GameManager, activeFilter *bool ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch all games from the manager
        log.Println("🎯 fetchGames called") 
		games := manager.GetAllGames()


		var response = []game.GameSnapshot{}

		for _, g := range games {
			if activeFilter == nil || g.Active == *activeFilter {
				response = append(response, g.Snapshot())
			}
		}

		// Set content type and return the JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
