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

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(game *g.Game, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    var player *g.Player

    for {
        messageType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            if player != nil {
                game.RemovePlayer(player)
            }
            return
        }

        fmt.Printf("Received: %s, Type: %d\n", msg, messageType)

        var m t.Message
        if err := json.Unmarshal(msg, &m); err != nil {
            log.Println("JSON unmarshal error:", err)
            continue
        }

        if m.Type == "join" && player == nil {
            player = game.AddPlayer(conn, m.Name)
            log.Printf("Player '%s'[#%s] joined the game\n", player.Name, player.ID)
            continue
        }

        if player == nil {
            log.Println("Message received before join; ignoring")
            continue
        }

        game.HandleMessage(player, m)
    }
}

