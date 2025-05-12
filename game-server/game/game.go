//game/game.go package game
package game

import (
    "encoding/json"
	  "fmt"
	  "log"
	  "math/rand"
	  "sync"
	  "time"
    "flip-cup/types" 
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
	  //Lobby     	    []*Player
	  TeamA     	    *Team
	  TeamB     	    *Team
	  QuestionFile  	*types.QuestionFile
	  Active  	      bool
	  mu        	    sync.Mutex
}

func NewGame(questionFile *types.QuestionFile) *Game {
    return &Game{
        //Lobby: []*Player{},
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

     // Log Teams
     fmt.Printf("📸 Game Snapshot: %+v\n", g.Snapshot())

    // Broadcast to all players
    g.Broadcast(types.Message{
        Type: "teams_update",
        Answer: g.Snapshot(),
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



func (g *Game) GetPlayer(ID string) *Player {
    if player := g.TeamA.GetPlayer(ID); player != nil {
        return player
    }
    if player := g.TeamB.GetPlayer(ID); player != nil {
        return player
    }
    return nil
}



func (g *Game) AddPlayer(conn *websocket.Conn, name string) *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := NewPlayer(conn, name)
	
  //g.Lobby = append(g.Lobby, player)
  if len(g.TeamB.Players) <  len(g.TeamA.Players) {
      g.TeamB.AddPlayer(player)
  } else {
      g.TeamA.AddPlayer(player)
  }
    
  g.Broadcast(types.Message{Type: "player_joined", Name: player.Name})
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


    g.Broadcast(types.Message{Type: "game_started"})

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

// empty teams into lobby
// shuffle lobby
// devide lobby players into teams
func (g *Game) ReassignTeams() {
	g.mu.Lock()
	defer g.mu.Unlock()


	if g.Active == true {
		fmt.Println("🚫 Cannot reassign teams in the middle of a game.")
	 	return
  }
    
   // Move players from both teams back into the temp lobby
  tmp := []*Player{}

  tmp = append(tmp, g.TeamA.Players...)
  tmp = append(tmp, g.TeamB.Players...)

  g.TeamA.Players = []*Player{}
  g.TeamB.Players = []*Player{}

	shufflePlayers(tmp)
/*
	if len(g.Lobby) % 2 != 0 {
	 	// its odd, we need an even team someone needs to leave
		fmt.Println("🚫 Cannot proceed: uneven number of players in the lobby.")
	 	return
 	}	
*/

	i := 0 
	for _, player := range tmp{
		if i % 2 == 0 {
			g.TeamA.Players = append(g.TeamA.Players, player)
		} else {
			g.TeamB.Players = append(g.TeamB.Players, player)
		}
		i++
	}
	tmp = []*Player{}
}

func (g *Game) HandleMessage(p *Player, m types.Message) {

	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println("Handle Message:\n", string(b))
	fmt.Printf("Message Type [%s]:\n", m.Type)

	team := g.GetTeam(p)

	switch m.Type {
	case "assign_teams":
		g.ReassignTeams()
		fallthrough
	case "show_players":
    g.DisplayGameSnapshot()
	case "join":
		p.Name = m.Name
		g.Broadcast(types.Message{Type: "player_joined", Name: p.Name})
		g.PlayerBroadcast(types.Message{Type: "joined_success", Name: p.Name}, p)
	case "answer":
		if g.Active == false {
		    fmt.Println("🚫 Cannot check answer for inactive games.")
    } else if g.CheckAnswer(p, team, m) {
			g.PlayerBroadcast(types.Message{Type: "correct", Name: p.Name}, p)
			if false ==  g.NextQuestion(team) {
        g.EndGame(team)
      }
      g.DisplayGameSnapshot() 
		} else {
			g.TeamBroadcast(types.Message{Type: "incorrect_answer", Name: p.Name}, team)
    }
	case "start":
    g.StartGame()
    g.DisplayGameSnapshot()
  //case "restart":
  //  g.RestartGame()
  //  g.DisplayGameSnapshot()
	}
}

func (g *Game) TeamBroadcast(msg types.Message, t *Team) {
	  for _, p := range t.Players {
        g.PlayerBroadcast(msg, p)
    }
}

func (g *Game) PlayerBroadcast(msg types.Message, p *Player) {
    data, _ := json.Marshal(msg)
    fmt.Print(msg)
		p.Conn.WriteMessage(websocket.TextMessage, data)
}

func (g *Game) Broadcast(msg types.Message) {
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
    g.Broadcast(types.Message{Type: "winner", Name: t.Name})
}

func (g *Game) NextQuestion(t *Team) bool{
	  if t.Turn > (len(t.Players) - 1) {
	  	  return false
	  }
    currentPlayer := t.GetCurrentPlayer()
	  q := g.QuestionFile.Questions[t.Turn]

    g.PlayerBroadcast(types.Message{Type: "question", Name: q.Prompt}, currentPlayer)
    return true
}

func (g *Game) CheckAnswer(p *Player, t *Team, answer types.Message) bool {

	  if ! t.IsPlayerAllowedToAnswer(p) {
        fmt.Println("wrong player answered")
        return false // the wrong player attempted to answer
	  }

    currentQuestion := g.QuestionFile.Questions[t.Turn] 
    if answer.Answer == currentQuestion.Answer {
        fmt.Println("correct answer")
	      t.Turn++
        return true
    }
    fmt.Println("incorrect answer")
    return false
}

func shufflePlayers(players []*Player) {
    rand.Seed(time.Now().UnixNano()) // Seed only once
    rand.Shuffle(len(players), func(i, j int) {
        players[i], players[j] = players[j], players[i]
    })
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func RandID() string {
	return fmt.Sprintf("%x", rand.Intn(999999))
}

func mustJSON(v any) string {
    b, err := json.Marshal(v)
    if err != nil {
        log.Fatalf("failed to marshal broadcast payload: %v", err)
    }
    return string(b)
}
