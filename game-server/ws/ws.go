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

    // Read initial join message first to get name before calling AddPlayer
    _, msg, err := conn.ReadMessage()
    if err != nil {
        log.Println("Read error during join:", err)
        return
    }

    var joinMsg t.Message
    if err := json.Unmarshal(msg, &joinMsg); err != nil {
        log.Println("JSON unmarshal error on join:", err)
        return
    }

    if joinMsg.Type != "join" {
        log.Println("Expected join message, got:", joinMsg.Type)
        return
    }

    player := game.AddPlayer(conn, joinMsg.Name)
    log.Printf("Player '%s'[#%s] joined the game\n", player.Name, player.ID)

    for {
        messageType, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Read error:", err)
            game.RemovePlayer(player)
            return
        }

        fmt.Printf("Received: %s, Type: %d\n", msg, messageType)

        var m t.Message
        if err := json.Unmarshal(msg, &m); err != nil {
            log.Println("JSON unmarshal error:", err)
            continue
        }

        game.HandleMessage(player, m)
    }
}
