// ws/handler.go
package ws

import (
	//"encoding/json"
	//"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"flip-cup/internal/game"
	"flip-cup/internal/quiz"
	"flip-cup/internal/transport/types"
	"flip-cup/internal/utils"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocketConnection(manager *game.GameManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
        log.Println("Incoming WS connection")

		conn, err := Upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading first message:", err)
			return
		}

		var envelope types.Envelope
		utils.MustUnmarshal(conn, msg, &envelope)

		utils.MustWriteJSON(conn, &envelope)
	
		// step 1, either create or join game
		// step 2, add player
		// step 3, join game loop

		var g *game.Game
		var p *game.Player
		var name string


		switch envelope.Type {
			case "join_existing_game":
				var joinPayload game.JoinExistingGamePayload
				utils.MustUnmarshal(conn, envelope.Payload, &joinPayload)

                log.Printf("Join Existing Game: %+v", envelope)
				
				// Handle joining an existing game
				g = manager.GetGame(joinPayload.GameID)
                g.DisplayGameSnapshot("joined_existing_game", p)

				if g == nil {
					log.Println("Error joining existing game:", err)
					return
				}

			case "create_game":
				// It's a Start Game payload
				var startPayload game.StartGamePayload
				utils.MustUnmarshal(conn, envelope.Payload, &startPayload)
  			    utils.LogPrettyJSON("Start New Game", envelope)
		
				var qf *quiz.QuestionFile
				qf, err := quiz.NewQuestionFile(startPayload.QuizFilename)
				if err != nil {
					log.Println("Error reading quiz file:", err)
					return
				}

				g = manager.NewGame(qf)
				log.Println("game_created: ", g.ID)
                g.DisplayGameSnapshot("ignore_created", p)

				conn.WriteJSON(
					types.Envelope{
						Type: "game_created",
						GameID: g.ID,
					},
				)

			default:
				log.Println("Unknown message type:", envelope.Type)
				return
		}
		p = g.AddPlayer(conn, name)

		// Handle the WebSocket messages for this game
		g.HandleConnection(p)
	}
}
