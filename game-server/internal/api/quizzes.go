// internal/api/quizzes.go
package api

import (
	"encoding/json"
	"fmt"
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

func fetchQuestionFiles() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

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
