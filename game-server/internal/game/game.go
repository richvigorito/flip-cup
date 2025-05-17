// game/game.go package game
package game

import (
	"encoding/json"
	"flip-cup/internal/transport/types"
	"flip-cup/internal/utils"
	"flip-cup/internal/quiz"
	"fmt"
	"time"
	"math/rand"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type GameSnapshot struct {
    ID    string            `json:"id"`
	  TeamA TeamSnapshot      `json:"teamA"`
	  TeamB TeamSnapshot      `json:"teamB"`
    Action    string            `json:"action,omitempty"'`
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
func (g *Game) DisplayTeamSnapshots() {
    snapshotTeamA := g.TeamA.Snapshot() 
    snapshotTeamB := g.TeamB.Snapshot() 

    aBytes, _ := json.Marshal(snapshotTeamA)
    g.TeamBroadcast(types.Envelope{
        Type: "my_current_team",
        Payload: aBytes,
    }, g.TeamA)

    bBytes, _ := json.Marshal(snapshotTeamB)
    g.TeamBroadcast(types.Envelope{
        Type: "my_current_team",
        Payload: bBytes,
    }, g.TeamB)
}

func (g *Game) DisplayGameSnapshot(action string, p *Player) {
    log.Printf("Game Active: %v\n", g.Active)

    var performedBy interface{}
    if p != nil {
        performedBy = p.Snapshot()
    } else {
        performedBy = nil 
    }

    snapshot := map[string]interface{}{
        "game_snapshot":        g.Snapshot(),
        "action_performed_by":  performedBy,
    }

    log.Printf("Game Snapshot: %+v\n", snapshot)
   
    snapshotBytes, _ := json.Marshal(snapshot)
    g.Broadcast(types.Envelope{
        Type: action,
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

func (g *Game) AddPlayer(conn *websocket.Conn, name string) *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := NewPlayer(conn, name)
	
  if len(g.TeamB.Players) <  len(g.TeamA.Players) {
		 log.Println("%s to team A: ", name)
      g.TeamB.AddPlayer(player)
  } else {
		 log.Println("%s to team B: ", name)
      g.TeamA.AddPlayer(player)
  }
    
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
    g.QuestionFile.ShuffleQuestions()


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
    g.QuestionFile.ShuffleQuestions()
    g.Active = false
}

func (g *Game) handleReassignTeams() {
    g.mu.Lock()
    defer g.mu.Unlock()

    if g.Active == true {
        log.Println("🚫 Cannot reassign teams in the middle of a game.")
        return
    }
        
    // Move players from both teams back into temp slice
    tmp := []*Player{}

    tmp = append(tmp, g.TeamA.Players...)
    tmp = append(tmp, g.TeamB.Players...)

    rand.Shuffle(len(tmp), func(i, j int) {
        tmp[i], tmp[j] = tmp[j],  tmp[i]
    })

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
        log.Println("🚫 Cannot check answer for inactive games.")
        return //  cannot answer active games
    }

    t := g.GetTeam(p)
    if ! t.IsPlayerAllowedToAnswer(p) {
        log.Println("wrong player answered")
        return // the wrong player attempted to answer
	  }

    currentQuestion := g.QuestionFile.Questions[t.Turn] 
    if answer.Answer == currentQuestion.Answer {
        log.Println("correct answer")
	      t.Turn++
        if false ==  g.NextQuestion(t) {
            g.EndGame(t)
        }
        g.DisplayGameSnapshot("answered_correctly", p)
    } else {
			g.TeamBroadcast(types.Envelope{Type: "incorrect_answer", Name: p.Name}, t)
    }
}

func (g *Game) handleAssignPlayerName(p *Player, addPlayerPayload *AddPlayerPayload) {
		p.Name = addPlayerPayload.Name
		g.Broadcast(types.Envelope{Type: "player_joined", Name: p.Name})
		g.PlayerBroadcast(types.Envelope{Type: "joined_success", Name: p.Name}, p)
}

//func (g *Game) HandleMessage(conn *websocket.Conn, msg types.Envelope) {
func (g *Game) HandleConnection(player *Player) {

   // player has entered game loop
   // boardcast they joined then 
   // start reading messages
    payload := PlayerJoinedPayload{
        PlayerID: player.ID,
    }
    g.PlayerBroadcast(types.Envelope{
        Type: "game_player_initialized",
        GameID: g.ID,
        Payload: utils.MustMarshal(payload),
    }, player)

    for{     
        _, msg, err := player.Conn.ReadMessage()
        if err != nil {
			      log.Println("Error reading message:", err)
			      return
		    }

        var envelope types.Envelope
        utils.MustUnmarshal(player.Conn, msg, &envelope)
    

        utils.LogPrettyJSON("🟢✉️↙️ Handle Envelope: ", envelope)

        g.HandleMessage(player, envelope)
    }
}

func (g *Game) HandleMessage(p *Player, msg types.Envelope) {

        switch msg.Type {

            case "add_player":
                var addPlayerPayload AddPlayerPayload
                utils.MustUnmarshal(p.Conn, msg.Payload, &addPlayerPayload)
                g.handleAssignPlayerName(p, &addPlayerPayload)
            case "assign_teams":
                g.handleReassignTeams()
                g.DisplayTeamSnapshots()
                time.Sleep(50 * time.Millisecond)
                g.DisplayGameSnapshot("teams_assigned", p)
            case "reassign_teams":
                g.handleReassignTeams()
                g.DisplayTeamSnapshots()
                time.Sleep(50 * time.Millisecond)
                g.DisplayGameSnapshot("teams_reassigned", p)
            case "show_players":
                g.DisplayGameSnapshot("show_players", p)
            case "check_answer":
                var answerPayload AnswerPayload
                utils.MustUnmarshal(p.Conn, msg.Payload, &answerPayload)
                g.handleCheckAnswer(p,  &answerPayload)
            case "start":
                g.StartGame()
                g.DisplayGameSnapshot("game_started", p)
          case "restart_game":
              g.RestartGame()
              g.DisplayGameSnapshot("game_restarted", p)
        }
}

func (g *Game) TeamBroadcast(msg types.Envelope, t *Team) {
	  for _, p := range t.Players {
        g.PlayerBroadcast(msg, p)
    }
}

func (g *Game) PlayerBroadcast(msg types.Envelope, p *Player) {
    data, _ := json.Marshal(msg)
    var logString = "🟦🔉↗️  Broadcast following Message to "+p.Name
    utils.LogPrettyJSON(logString, msg)
		p.Conn.WriteMessage(websocket.TextMessage, data)
}

func (g *Game) Broadcast(msg types.Envelope) {
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

func RandID() string {
	return fmt.Sprintf("%x", rand.Intn(999999))
}
