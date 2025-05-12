// http/http.go
package http

import (
	"fmt"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"github.com/gorilla/mux"
  "gopkg.in/yaml.v2"
	"io/ioutil"
  "path/filepath"
	"flip-cup/game"
)

// SetupRoutes sets up HTTP endpoints
func SetupRoutes(manager *game.GameManager, r *mux.Router) {
	r.HandleFunc("/games", fetchGames(manager))
	r.HandleFunc("/question-files", fetchQuestionFiles())
}

func fetchGames(manager *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch all games from the manager
		games := manager.GetAllGames()

		var response = []game.GameSnapshot{}

		for _, g := range games {
			if (!g.Active){
				response = append(response, g.Snapshot())
			}
		}

		// Set content type and return the JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}


type QuestionMeta struct {
    File  string `json:"file"`
    Label string `json:"label"`
    Category string `json:"category"`
}

func fetchQuestionFiles() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        files, err := os.ReadDir("questions")
        if err != nil {
            http.Error(w, "Failed to read questions directory", http.StatusInternalServerError)
            return
        }

        var result []QuestionMeta

        for _, f := range files {
            if !f.IsDir() && strings.HasSuffix(f.Name(), ".yaml") {
                fullPath := filepath.Join("questions", f.Name())
                content, err := ioutil.ReadFile(fullPath)
                if err != nil {
                    fmt.Println("skip unreadable file: ", fullPath)
                    continue // skip unreadable file
                }

                var parsed map[string]interface{}
                if err := yaml.Unmarshal(content, &parsed); err != nil {
                    fmt.Println("skip unparsable YAML: ", fullPath)
                    continue // skip unparsable YAML
                }

                title := f.Name() // fallback if no title found
                category := ""
                if parsedTitle, ok := parsed["title"].(string); ok {
                    title = parsedTitle
                }
                if parsedCategory, ok := parsed["category"].(string); ok {
                    category = parsedCategory
                }

                result = append(result, QuestionMeta{
                    File:  f.Name(),
                    Label: title,
                    Category: category,
                })
            }
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}
