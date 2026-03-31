// internal/transport/api/router.go
package api

import (
	"time"

	"flip-cup/internal/game"
	"github.com/gorilla/mux"
)

// SetupRoutes sets up HTTP endpoints.
func SetupRoutes(manager *game.GameManager, staleAfter time.Duration, r *mux.Router) {
	r.HandleFunc("/games/{status:active|inactive|stale}", fetchGames(manager, staleAfter)).Methods("GET")
	r.HandleFunc("/games/stale", deleteStaleGames(manager, staleAfter)).Methods("DELETE")
	r.HandleFunc("/games/{id}", deleteGame(manager)).Methods("DELETE")
	r.HandleFunc("/quizzes", fetchQuestionFiles())
	r.PathPrefix("/").Handler(SPAHandler("public"))
}
