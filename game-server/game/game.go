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
	  "io/ioutil"
    "gopkg.in/yaml.v2"
)

type QuestionFile struct {
    Questions []*types.Question `yaml:"questions"`
}

type Game struct {
	  Lobby     	[]*Player
	  TeamA     	*Team
	  TeamB     	*Team
	  Questions  	[]*types.Question
	  Active  	  bool
	  mu        	sync.Mutex
}

func (g *Game) LoadQuestions(filename string) error {
    data, err := ioutil.ReadFile(filename) // Read the YAML file
    if err != nil {
        return fmt.Errorf("error reading YAML file: %v", err)
    }

    var questions []*types.Question  // Make sure this matches the expected structure
    err = yaml.Unmarshal(data, &questions)
    if err != nil {
        return fmt.Errorf("error unmarshaling YAML: %v", err)
    }

    g.Questions = questions  // Assign questions to the Game struct
    fmt.Println("Questions Loaded:")
    for _,q := range g.Questions {
	    fmt.Println("-",q.Prompt)
    }
    return nil
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



func NewGame() *Game {
    return &Game{
        Lobby: []*Player{},
        TeamA: &Team{Players: []*Player{}, Name: "A-Team",  Turn: 0},
        TeamB: &Team{Players: []*Player{}, Name: "B-squad", Turn: 0},
        Questions: []*types.Question{},
        Active: false,
    }
}

func (g *Game) AddPlayer(conn *websocket.Conn, name string) *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := &Player{
		ID:   RandID(),
		Conn: conn,
		Name: name,
	//	Index: len(g.Lobby),
	}
	g.Lobby = append(g.Lobby, player)
  g.Broadcast(types.Message{Type: "player_joined", Name: player.Name})
	return player
}

func (g *Game) RemovePlayer(p *Player) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i, pl := range g.Lobby {
		if pl == p {
			g.Lobby = append(g.Lobby[:i], g.Lobby[i+1:]...)
			break
		}
	}
}

func (g *Game) StartGame() {
    if len(g.Questions) == 0 {
        fmt.Println("No questions loaded.")
        return
    }

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
    if len(g.Questions) == 0 {
        fmt.Println("No questions loaded.")
        return
    }

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

   	// Move players from both teams back into the lobby
    	g.Lobby = append(g.Lobby, g.TeamA.Players...)
    	g.Lobby = append(g.Lobby, g.TeamB.Players...)

    	g.TeamA.Players = []*Player{}
    	g.TeamB.Players = []*Player{}

	shufflePlayers(g.Lobby)

	if len(g.Lobby) % 2 != 0 {
	 	// its odd, we need an even team someone needs to leave
		fmt.Println("🚫 Cannot proceed: uneven number of players in the lobby.")
	 	return
 	}	

	i := 0 
	for _, player := range g.Lobby{
		if i % 2 == 0 {
			g.TeamA.Players = append(g.TeamA.Players, player)
		} else {
			g.TeamB.Players = append(g.TeamB.Players, player)
		}
		i++
	}
	g.Lobby = []*Player{}
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
    fmt.Printf("Game Active: %v\n", g.Active)

    // Log Lobby
    fmt.Println("Lobby:")
    for _, player := range g.Lobby {
        fmt.Printf("  '%s' [#%s]\n", player.Name, player.ID)
    }

    // Get team snapshots
    snapA := g.TeamA.snapshot()
    snapB := g.TeamB.snapshot()

    // Log Teams
    fmt.Printf("TeamA (Turn %d): %v\n", snapA.Turn, snapA.Players)
    fmt.Printf("TeamB (Turn %d): %v\n", snapB.Turn, snapB.Players)

    // Broadcast to all players
    g.Broadcast(types.Message{
        Type: "teams_update",
        Answer: map[string]any{
            "teamA": g.TeamA.snapshot(),
            "teamB": g.TeamB.snapshot(),
        },
    }) 
	case "join":
		p.Name = m.Name
		g.Broadcast(types.Message{Type: "player_joined", Name: p.Name})
	case "answer":
		if g.Active == false {
		    fmt.Println("🚫 Cannot check answer for inactive games.")
    } else if g.CheckAnswer(p, team, m) {
			g.Broadcast(types.Message{Type: "correct", Name: p.Name})
			if false ==  g.NextQuestion(team) {
        g.EndGame(team)
      }
		} else {
			g.TeamBroadcast(types.Message{Type: "incorrect_answer", Name: p.Name}, team)
    }
	case "start":
        	g.StartGame()
  case "restart":
        	g.RestartGame()
	}
}

func (g *Game) TeamBroadcast(msg types.Message, t *Team) {
	  for _, p := range t.Players {
        g.PlayerBroadcast(msg, p)
    }
}

func (g *Game) PlayerBroadcast(msg types.Message, p *Player) {
    data, _ := json.Marshal(msg)
		p.Conn.WriteMessage(websocket.TextMessage, data)
}

func (g *Game) Broadcast(msg types.Message) {
	data, _ := json.Marshal(msg)
	for _, p := range g.Lobby {
		p.Conn.WriteMessage(websocket.TextMessage, data)
	}
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
	  q := g.Questions[t.Turn]
    //p = t.Players[t
    //TODO LEFT OFF HERE
    g.Broadcast(types.Message{Type: "question", Name: q.Prompt})
    return true
}

func (g *Game) CheckAnswer(p *Player, t *Team, answer types.Message) bool {

	  if ! t.IsPlayerAllowedToAnswer(p) {
        fmt.Println("wrong player answered")
        return false // the wrong player attempted to answer
	  }

    currentQuestion := g.Questions[t.Turn] 
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
