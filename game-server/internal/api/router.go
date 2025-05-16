// internal/api/router.go
package api

import (
	"github.com/gorilla/mux"
	"flip-cup/internal/game"
)

// SetupRoutes sets up HTTP endpoints
func SetupRoutes(manager *game.GameManager, r *mux.Router) {
	inactive := false /* ie get inactive games only */
	r.HandleFunc("/games/active", fetchGames(manager, &inactive))
	r.HandleFunc("/quizzes", fetchQuestionFiles())
}

