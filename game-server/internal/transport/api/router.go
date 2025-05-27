// internal/transport/api/router.go
package api

import (
	"github.com/gorilla/mux"
	"flip-cup/internal/game"
)

// SetupRoutes sets up HTTP endpoints 
func SetupRoutes(manager *game.GameManager, r *mux.Router) {
	r.HandleFunc("/games/{status}", fetchGames(manager)).Methods("GET")
	r.HandleFunc("/games/{status}", deleteGames(manager)).Methods("DELETE")
	r.HandleFunc("/games/{id}", deleteGame(manager)).Methods("DELETE")
	r.HandleFunc("/quizzes", fetchQuestionFiles())
    r.PathPrefix("/").Handler(SPAHandler("public"))

/**
	$active ... get from url 
	r.HandleFunc("/games/active", fetchGames(manager, $active))
*/

}

