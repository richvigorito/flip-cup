package utils

import (
	"encoding/json"
	"log"
)

func LogPrettyJSON(prefix string, data any) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("%s: error marshaling JSON: %v", prefix, err)
		return
	}
	log.Printf("%s:\n%s", prefix, string(b))
}

