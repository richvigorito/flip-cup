// game/game.go package game
package game

import (
	"encoding/json"
	"flip-cup/internal/transport/types"
	"flip-cup/internal/utils"
	"flip-cup/internal/quiz"
	"fmt"
	"math/rand"
	//"log"
	"sync"

	"github.com/gorilla/websocket"
)

type GameSnapshot struct {
    ID    string            `json:"id"`
	  TeamA TeamSnapshot      `json:"teamA"`
	  TeamB TeamSnapshot      `json:"teamB"`
    Active bool             `json:"active"`
}

type Game struct {
	  ID              string
	  TeamA     	    *Team
	  TeamB     	    *Team
	  QuestionFile  	*quiz.QuestionFile
	  Active  	      bool
	  mu        	    sync.Mutex
}

func NewGame(questionFile *quiz.QuestionFile) *Game {
    return &Game{
        TeamA: &Team{Players: []*Player{}, Name: "A-Team",  Turn: 0},
        TeamB: &Team{Players: []*Player{}, Name: "B-squad", Turn: 0},
        QuestionFile: questionFile,
        Active: false,
    }
}

func (g *Game) Snapshot() GameSnapshot {
    return GameSnapshot{ 
        ID: g.ID,
        TeamA: g.TeamA.Snapshot(),
        TeamB: g.TeamB.Snapshot(),
        Active: g.Active,
    }
}

func (g *Game) DisplayGameSnapshot() {
    fmt.Printf("Game Active: %v\n", g.Active)

    fmt.Printf("Game Snapshot: %+v\n", g.Snapshot())
   
    snapshotBytes, _ := json.Marshal(g.Snapshot())
    g.Broadcast(types.Envelope{
        Type: "game_snapshot",
        Payload: snapshotBytes,
    })
}

func (g *Game) GetTeam(p *Player) *Team {
    for _, player := range g.TeamB.Players {
        if p.ID == player.ID {
            return g.TeamB
        }
    }

    for _, player := range g.TeamA.Players {
        if p.ID == player.ID {
            return g.TeamA
        }
    }

    return nil
}
/*
func (g *Game) GetPlayerByConn(conn *websocket.Conn) *Player {
	for _, p := range g.TeamA.Players {
		if p.Conn == conn {
			return p
		}
	}
	for _, p := range g.TeamB.Players {
		if p.Conn == conn {
			return p
		}
	}
	return nil
}
*/


/*
func (g *Game) GetPlayer(ID string) *Player {
    if player := g.TeamA.GetPlayer(ID); player != nil {
        return player
    }
    if player := g.TeamB.GetPlayer(ID); player != nil {
        return player
    }
    return nil
}
*/



//func (g *Game) handleAddPlayer(conn *websocket.Conn, payload *AddPlayerPayload) *Player {
func (g *Game) AddPlayer(conn *websocket.Conn, name string) *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	//player := NewPlayer(conn, payload.Name)
	player := NewPlayer(conn, name)
	
  if len(g.TeamB.Players) <  len(g.TeamA.Players) {
      g.TeamB.AddPlayer(player)
  } else {
      g.TeamA.AddPlayer(player)
  }
    
  //g.Broadcast(types.Envelope{Type: "player_joined", Name: player.Name})
	return player
}

func (g *Game) RemovePlayer(p *Player) {
	g.mu.Lock()
	defer g.mu.Unlock()
  g.TeamA.RemovePlayer(p)
  g.TeamB.RemovePlayer(p)
}


func (g *Game) StartGame() {
    // reset everything
    g.Active = true
    g.TeamA.Turn = 0 
    g.TeamB.Turn = 0 


    g.Broadcast(types.Envelope{Type: "game_started"})

    // Start first question for both teams
    if len(g.TeamA.Players) > 0 {
        g.NextQuestion(g.TeamA)
    }
    if len(g.TeamB.Players) > 0 {
        g.NextQuestion(g.TeamB)
    }
}

func (g *Game) RestartGame() {
    // Start first question for both teams
    g.TeamA.Turn = 0
    g.TeamB.Turn = 0
}

func (g *Game) handleReassignTeams() {
    g.mu.Lock()
    defer g.mu.Unlock()

    if g.Active == true {
        fmt.Println("🚫 Cannot reassign teams in the middle of a game.")
        return
    }
        
    // Move players from both teams back into temp slice
    tmp := []*Player{}

    tmp = append(tmp, g.TeamA.Players...)
    tmp = append(tmp, g.TeamB.Players...)

    g.TeamA.Players = []*Player{}
    g.TeamB.Players = []*Player{}

    i := 0 
    for _, player := range tmp{
        if i % 2 == 0 {
            g.TeamA.Players = append(g.TeamA.Players, player)
        } else {
          g.TeamB.Players = append(g.TeamB.Players, player)
        }
        i++
    }
    g.TeamA.Shuffle()
    g.TeamB.Shuffle()
    tmp = []*Player{}
}

func (g *Game) handleCheckAnswer(p *Player, answer *AnswerPayload) {
    if g.Active == false {
        fmt.Println("🚫 Cannot check answer for inactive games.")
        return //  cannot answer active games
    }

    t := g.GetTeam(p)
    if ! t.IsPlayerAllowedToAnswer(p) {
        fmt.Println("wrong player answered")
        return // the wrong player attempted to answer
	  }

    currentQuestion := g.QuestionFile.Questions[t.Turn] 
    if answer.Answer == currentQuestion.Answer {
        fmt.Println("correct answer")
	      t.Turn++
        if false ==  g.NextQuestion(t) {
            g.EndGame(t)
        }
        g.DisplayGameSnapshot()
    } else {
			g.TeamBroadcast(types.Envelope{Type: "incorrect_answer", Name: p.Name}, t)
    }
}
/*
func (g *Game) handleAssignPlayerName(p *Player, answer *AnswerPayload) {
		p.Name = name
		g.Broadcast(types.Envelope{Type: "player_joined", Name: p.Name})
		g.PlayerBroadcast(types.Envelope{Type: "joined_success", Name: p.Name}, p)
}
*/

//func (g *Game) HandleMessage(conn *websocket.Conn, msg types.Envelope) {
func (g *Game) HandleConnection(player *Player) {

   //  player has entered game loop
   //  boardcast they joined then 
   // start reading messages
   payload := PlayerJoinedPayload{
        PlayerID: player.ID,
        Name: player.Name,
    }
    g.Broadcast(types.Envelope{
        Type: "player_joied",
        GameID: g.ID,
        Payload: utils.MustMarshal(payload),
    })

    for{     
        _, msg, err := player.Conn.ReadMessage()
        if err != nil {
			      fmt.Println("Error reading message:", err)
			      //log.Println("Error reading message:", err)
			      return
		    }

        var envelope types.Envelope
        utils.MustUnmarshal(player.Conn, msg, &envelope)
    

        fmt.Printf("Handle Envelope: %+v", envelope)
        fmt.Printf("Message Type [%s]:\n", envelope.Type)

        g.HandleMessage(player, envelope)
    }
}

func (g *Game) HandleMessage(p *Player, msg types.Envelope) {

        switch msg.Type {
            case "assign_teams":
                g.handleReassignTeams()
                fallthrough
            case "show_players":
                g.DisplayGameSnapshot()
            //case "add_player":
                //var addPlayerPayload AddPlayerPayload
                //utils.MustUnmarshal(conn, envelope.Payload, &addPlayerPayload)
                //g.handleAddPlayer(conn, &addPlayerPayload)
            case "check_answer":
                var answerPayload AnswerPayload
                utils.MustUnmarshal(p.Conn, msg.Payload, &answerPayload)
                //p := g.GetPlayerByConn(conn)
                //if p == nil {
                    //fmt.Println("Could not find player for connection")
                //return
            //}  
                g.handleCheckAnswer(p,  &answerPayload)
            case "start":
                g.StartGame()
                g.DisplayGameSnapshot()
          //case "restart":
            //  g.RestartGame()
            //  g.DisplayGameSnapshot()
        }
}

func (g *Game) TeamBroadcast(msg types.Envelope, t *Team) {
	  for _, p := range t.Players {
        g.PlayerBroadcast(msg, p)
    }
}

func (g *Game) PlayerBroadcast(msg types.Envelope, p *Player) {
    data, _ := json.Marshal(msg)
    fmt.Print(msg)
		p.Conn.WriteMessage(websocket.TextMessage, data)
}

func (g *Game) Broadcast(msg types.Envelope) {
	/*data, _ := json.Marshal(msg)
	for _, p := range g.Lobby {
		p.Conn.WriteMessage(websocket.TextMessage, data)
	}
  */
  fmt.Printf("Broadcast following Message: %+v\n", msg)
  g.TeamBroadcast(msg, g.TeamA)
  g.TeamBroadcast(msg, g.TeamB)
}


func (g *Game) EndGame(t *Team) {
    g.Active = false
    g.Broadcast(types.Envelope{Type: "winner", Name: t.Name})
}

func (g *Game) NextQuestion(t *Team) bool{
	  if t.Turn > (len(t.Players) - 1) {
	  	  return false
	  }
    currentPlayer := t.GetCurrentPlayer()
	  q := g.QuestionFile.Questions[t.Turn]

    g.PlayerBroadcast(types.Envelope{Type: "question", Name: q.Prompt}, currentPlayer)
    return true
}

/*
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
*/
func RandID() string {
	return fmt.Sprintf("%x", rand.Intn(999999))
}
/*
func mustJSON(v any) string {
    b, err := json.Marshal(v)
    if err != nil {
        log.Fatalf("failed to marshal broadcast payload: %v", err)
    }
    return string(b)
}
*/
