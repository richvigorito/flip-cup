// ws/ws.go
package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	g "flip-cup/game"
	t "flip-cup/types"
)

// Upgrader for WebSocket connection
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocket handles the WebSocket connection and game-specific logic
func HandleWebSocket(game *g.Game, conn *websocket.Conn) {
	var player *g.Player

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			if player != nil {
				game.TeamA.RemovePlayer(player)
				game.TeamB.RemovePlayer(player)
			}
			return
		}

		fmt.Printf("Received: %s, Type: %d\n", msg, messageType)

		var m t.Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		// Handle player join
		if m.Type == "join" && player == nil {
			player = game.AddPlayer(conn, m.Name)
			log.Printf("Player '%s'[#%s] joined the game\n", player.Name, player.ID)
			continue
		}

		// Ensure player is joined before handling other messages
		if player == nil {
			log.Println("Message received before join; ignoring")
			continue
		}

		// Handle other game-specific messages
		game.HandleMessage(player, m)
	}
}
