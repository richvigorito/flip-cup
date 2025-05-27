// internal/transport/api/quizzes.go
package api

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"net/http"
)

type QuizMeta struct {
	File     string `json:"file"`
	Label    string `json:"label"`
	Category string `json:"category"`
}


//
// GET: /api/quizzies
//
func fetchQuestionFiles() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println("🎯 fetchQuestionFiles called") 

        files, err := os.ReadDir("questions")
        if err != nil {
            http.Error(w, "Failed to read questions directory", http.StatusInternalServerError)
            return
        }

        var result []QuizMeta

        for _, f := range files {
            if !f.IsDir() && strings.HasSuffix(f.Name(), ".yaml") {
                fullPath := filepath.Join("questions", f.Name())
                content, err := os.ReadFile(fullPath)
                if err != nil {
                    log.Println("skip unreadable file: ", fullPath)
                    continue // skip unreadable file
                }

                var parsed map[string]interface{}
                if err := yaml.Unmarshal(content, &parsed); err != nil {
                    log.Println("skip unparsable YAML: ", fullPath)
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

                result = append(result, QuizMeta{
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
