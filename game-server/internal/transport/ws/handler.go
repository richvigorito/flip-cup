// internal/transport/ws/handler.go
package ws

import (
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

		switch envelope.Type {
		case "join_existing_game":
			var joinPayload game.JoinExistingGamePayload
			utils.MustUnmarshal(conn, envelope.Payload, &joinPayload)

			log.Printf("Join Existing Game: %+v", envelope)

			// Handle joining an existing game
			g = manager.GetGame(joinPayload.GameID)
			if g == nil {
				log.Printf("Error joining existing game: game %s not found", joinPayload.GameID)
				return
			}

			conn.WriteJSON(types.Envelope{
				Type: "joined_existing_game",
				Payload: utils.MustMarshal(map[string]interface{}{
					"game_snapshot": g.Snapshot(),
				}),
			})

			if joinPayload.PlayerID != "" {
				var t *game.Team
				p, t = g.ReconnectPlayer(joinPayload.PlayerID, conn)
				if p != nil {
					log.Printf("Reconnected player %s to game %s", p.ID, g.ID)

					// Send me info FIRST
					payload := game.PlayerJoinedPayload{
						PlayerID: p.ID,
						Name:     p.Name,
					}
					g.PlayerBroadcast(types.Envelope{
						Type:    "game_player_initialized",
						GameID:  g.ID,
						Payload: utils.MustMarshal(payload),
					}, p)

					// Send current state to reconnected player
					snapshot := map[string]interface{}{
						"game_snapshot":       g.Snapshot(),
						"action_performed_by": p.Snapshot(),
					}

					// Always send team info on reconnect
					if t != nil {
						log.Printf("Restoring team state for player %s (Team: %s)", p.Name, t.Name)
						teamData := utils.MustMarshal(t.Snapshot())
						g.PlayerBroadcast(types.Envelope{
							Type:    "my_current_team",
							Payload: teamData,
						}, p)
					} else {
						log.Printf("⚠️ Reconnected player %s has no team!", p.ID)
					}

					if g.Active {
						g.PlayerBroadcast(types.Envelope{Type: "game_started", Payload: utils.MustMarshal(snapshot)}, p)

						// Resend current question
						if t != nil {
							if t.Turn < len(g.QuestionFile.Questions) {
								q := g.QuestionFile.Questions[t.Turn]
								g.PlayerBroadcast(types.Envelope{Type: "question", Name: q.Prompt}, p)
							} else {
								log.Printf("⚠️ Team %s turn %d exceeds question count %d", t.Name, t.Turn, len(g.QuestionFile.Questions))
							}
						}
					} else {
						// For inactive games, send teams_assigned to sync lobby state
						g.PlayerBroadcast(types.Envelope{Type: "teams_assigned", Payload: utils.MustMarshal(snapshot)}, p)
					}
				}
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
					Type:   "game_created",
					GameID: g.ID,
				},
			)

		default:
			log.Println("Unknown message type:", envelope.Type)
			return
		}

		// From here, game is responsible for handling game specific messages
		g.HandleConnection(conn, p)
	}
}
