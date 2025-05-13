// internal/utils/errors.go

package utils

import (
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)

// ErrorHandler 
func HandleError(conn *websocket.Conn, err error, msg string) {
	log.Printf("%s: %v\n", msg, err)

	// Optionally, you can send an error message to the client
	errorMessage := fmt.Sprintf("Error: %s", err.Error())
	if conn != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(errorMessage))
	}
}

