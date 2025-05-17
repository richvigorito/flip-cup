package quiz

import (
    "fmt"
    "regexp"
    "strings"
)

type Question struct {
    Prompt string `yaml:"prompt"`
    Answers []string `yaml:"answers"`
}

func (q *Question) CheckAnswer(input string) bool {
    normalizedInput := normalize(input)
    for _, correct := range q.Answers {
        fmt.Println("input: ", normalizedInput, "correct", normalize(correct),  "condtinional:", normalize(correct) == normalizedInput)
        if normalize(correct) == normalizedInput {
            return true
        }
    }
    return false
}

var nonAlphaNum = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func normalize(s string) string {
    s = strings.TrimSpace(s)
    s = strings.ToLower(s)
    s = nonAlphaNum.ReplaceAllString(s, "")
    return s
}
