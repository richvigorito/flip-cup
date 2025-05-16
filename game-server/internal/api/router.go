// internal/api/router.go
package api

import (
    "net/http"
    "path/filepath"
	"github.com/gorilla/mux"
	"flip-cup/internal/game"
)

// SetupRoutes sets up HTTP endpoints 
func SetupRoutes(manager *game.GameManager, r *mux.Router) {
	inactive := false /* ie get inactive games only */
	r.HandleFunc("/games/active", fetchGames(manager, &inactive))
	r.HandleFunc("/quizzes", fetchQuestionFiles())

    // Serves static files from ./public
    staticDir := http.Dir("public")
    staticFileHandler := http.StripPrefix("/", http.FileServer(staticDir))

    r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

        // ignore ws (websocket route from prefix) 
	    if req.URL.Path == "/ws" || req.URL.Path == "/ws-test" {
			http.NotFound(w, req)
			return
		}

        path := filepath.Join("public", req.URL.Path)
        _, err := staticDir.Open(path) /
        if err != nil {
            // If file doesn't exist, serve index.html (for Svelte SPA routing)
            http.ServeFile(w, req, "public/index.html")
            return
        }
        staticFileHandler.ServeHTTP(w, req)
    })
}

