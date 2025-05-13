// internal/utils/must.go
package utils

import (
	"encoding/json"
	//"fmt"
//	"log"
//  "flip-cup/internal/game"
	"github.com/gorilla/websocket"	
//	"net/http"

)

func MustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic("failed to marshal payload: " + err.Error())
	}
	return data
}


func MustUnmarshal(conn *websocket.Conn, payload []byte, v interface{}) {
	if err := json.Unmarshal(payload, v); err != nil {
		HandleError(conn, err, "MustUnmarshal failed:")
		panic(err)  
	}
}

func MustWriteJSON(conn *websocket.Conn, v interface{}) {
	if err := conn.WriteJSON(v); err != nil {
		HandleError(conn, err, "MustWriteJSON failed:")
		panic(err)  
	}
}

/*
// MustGetGame is a helper that retrieves a game and panics if it's not found.
func MustGetGame(manager *game.GameManager, gameID string) *game.Game {
	game := manager.GetGame(gameID)
	if game == nil {
		log.Fatalf("Game with ID %s not found", gameID)
	}
	return game
}

*/
