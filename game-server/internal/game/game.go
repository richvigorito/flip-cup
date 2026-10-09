package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"flip-cup/internal/quiz"
	"flip-cup/internal/transport/types"
	"flip-cup/internal/utils"
)

const (
	MinQuestionsPerPlayer = 1
	MaxQuestionsPerPlayer = 5
)

// SkipPenalty is added to the team's wait for every skip in a turn.
// main overrides it from FLIPCUP_SKIP_PENALTY (used by the E2E suite).
var SkipPenalty = 7 * time.Second

type GameSnapshot struct {
	ID                 string       `json:"id"`
	TeamA              TeamSnapshot `json:"teamA"`
	TeamB              TeamSnapshot `json:"teamB"`
	Action             string       `json:"action,omitempty"'`
	QuizFile           string       `json:"quizfile,omitempty"'`
	Active             bool         `json:"active"`
	QuestionsPerPlayer int          `json:"questionsPerPlayer"`
}

type Game struct {
	ID                 string
	TeamA              *Team
	TeamB              *Team
	QuestionFile       *quiz.QuestionFile
	Active             bool
	QuestionsPerPlayer int
	LastActivity       time.Time
	mu                 sync.Mutex

	// flow serializes turn/skip/penalty state; the penalty timer runs on another goroutine.
	flow       sync.Mutex
	generation int

	// Test seams: nil/default in production.
	afterFunc func(d time.Duration, f func()) (cancel func())
	onSend    func(p *Player, msg types.Envelope)
}

func NewGame(questionFile *quiz.QuestionFile) *Game {
	return &Game{
		ID:           utils.RandID(),
		TeamA:        &Team{Players: []*Player{}, Name: "A-Team", Turn: 0},
		TeamB:        &Team{Players: []*Player{}, Name: "B-squad", Turn: 0},
		QuestionFile: questionFile,
		Active:       false,
		QuestionsPerPlayer: MinQuestionsPerPlayer,
		LastActivity: time.Now(),
		afterFunc: func(d time.Duration, f func()) func() {
			t := time.AfterFunc(d, f)
			return func() { t.Stop() }
		},
	}
}

func (g *Game) UpdateActivity() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.LastActivity = time.Now()
}

// Snapshot & Display
func (g *Game) Snapshot() GameSnapshot {
	return GameSnapshot{
		ID:       g.ID,
		TeamA:    g.TeamA.Snapshot(),
		TeamB:    g.TeamB.Snapshot(),
		QuizFile: g.QuestionFile.Filename,
		Active:   g.Active,
		QuestionsPerPlayer: g.QuestionsPerPlayer,
	}
}

func (g *Game) DisplayTeamSnapshots() {
	for _, team := range []*Team{g.TeamA, g.TeamB} {
		data, _ := json.Marshal(team.Snapshot())
		g.TeamBroadcast(types.Envelope{
			Type:    "my_current_team",
			Payload: data,
		}, team)
	}
}

func (g *Game) DisplayGameSnapshot(action string, p *Player) {
	log.Printf("Game Active: %v\n", g.Active)
	snapshot := map[string]interface{}{
		"game_snapshot":       g.Snapshot(),
		"action_performed_by": safeSnapshot(p),
	}
	data, _ := json.Marshal(snapshot)

	g.Broadcast(types.Envelope{
		Type:    action,
		Payload: data,
	})
}

func safeSnapshot(p *Player) interface{} {
	if p != nil {
		return p.Snapshot()
	}
	return nil
}

// Broadcast
func (g *Game) Broadcast(msg types.Envelope) {
	g.TeamBroadcast(msg, g.TeamA)
	g.TeamBroadcast(msg, g.TeamB)
}

func (g *Game) TeamBroadcast(msg types.Envelope, t *Team) {
	for _, p := range t.Players {
		g.PlayerBroadcast(msg, p)
	}
}

func (g *Game) IsStale(threshold time.Duration) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return time.Since(g.LastActivity) > threshold
}

func (g *Game) PlayerBroadcast(msg types.Envelope, p *Player) {
	g.UpdateActivity()
	if g.onSend != nil {
		g.onSend(p, msg)
	}
	data, _ := json.Marshal(msg)
	var logString = "🟦🔉↗️  Broadcast following Message to " + p.Name
	utils.LogPrettyJSON(logString, msg)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Conn == nil {
		return
	}
	p.Conn.WriteMessage(websocket.TextMessage, data)
}

// Player Management
func (g *Game) ReconnectPlayer(playerID string, conn *websocket.Conn) (*Player, *Team) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, team := range []*Team{g.TeamA, g.TeamB} {
		for _, p := range team.Players {
			if p.ID == playerID {
				p.mu.Lock()
				p.Conn = conn
				p.mu.Unlock()
				return p, team
			}
		}
	}
	return nil, nil
}

func (g *Game) AddPlayer(conn *websocket.Conn, name string) *Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	player := NewPlayer(conn, name)
	if len(g.TeamB.Players) < len(g.TeamA.Players) {
		log.Printf("%s to team A: ", name)
		g.TeamB.AddPlayer(player)
	} else {
		log.Printf("%s to team B: ", name)
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

func (g *Game) GetTeam(p *Player) *Team {
	for _, team := range []*Team{g.TeamA, g.TeamB} {
		for _, player := range team.Players {
			if p.ID == player.ID {
				return team
			}
		}
	}
	return nil
}

// Game Flow
func (g *Game) StartGame(p *Player) {
	g.flow.Lock()
	defer g.flow.Unlock()

	g.Active = true
	g.generation++
	g.TeamA.ResetProgress()
	g.TeamB.ResetProgress()
	g.QuestionFile.ShuffleQuestions()

	// Broadcast game started BEFORE sending questions
	g.DisplayGameSnapshot("game_started", p)

	// Start first question for both teams
	if len(g.TeamA.Players) > 0 {
		g.NextQuestion(g.TeamA)
	}
	if len(g.TeamB.Players) > 0 {
		g.NextQuestion(g.TeamB)
	}
}

func (g *Game) RestartGame() {
	g.flow.Lock()
	defer g.flow.Unlock()

	g.generation++
	g.TeamA.ResetProgress()
	g.TeamB.ResetProgress()
	g.QuestionFile.ShuffleQuestions()
	g.Active = false
}

// SetQuestionsPerPlayer clamps count to the allowed range. It is a no-op once the game has started.
func (g *Game) SetQuestionsPerPlayer(count int) bool {
	g.flow.Lock()
	defer g.flow.Unlock()

	if g.Active {
		return false
	}
	if count < MinQuestionsPerPlayer {
		count = MinQuestionsPerPlayer
	}
	if count > MaxQuestionsPerPlayer {
		count = MaxQuestionsPerPlayer
	}
	g.QuestionsPerPlayer = count
	return true
}

func (g *Game) teamTotal(t *Team) int {
	return len(t.Players) * g.QuestionsPerPlayer
}

// questionAt wraps so skips can never run past the end of the quiz.
func (g *Game) questionAt(idx int) *quiz.Question {
	qs := g.QuestionFile.Questions
	return qs[idx%len(qs)]
}

// NextQuestion asks the team's current player their question. Callers hold g.flow.
func (g *Game) NextQuestion(t *Team) bool {
	if t.Turn >= g.teamTotal(t) {
		return false
	}
	currentPlayer := t.GetCurrentPlayer()
	q := g.questionAt(t.QuestionIdx)

	g.PlayerBroadcast(types.Envelope{Type: "question", Name: q.Prompt}, currentPlayer)
	return true
}

func (g *Game) EndGame(t *Team) {
	g.Active = false
	g.Broadcast(types.Envelope{Type: "winner", Name: t.Name})
}

func (g *Game) UpdateQuiz(p *Player, payload *UpdateQuiz) {
	qf, err := quiz.NewQuestionFile(payload.Filename)
	if err != nil {
		log.Println("Error reading quiz file:", err)
		return
	}
	g.QuestionFile = qf
}

// Handlers
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
		tmp[i], tmp[j] = tmp[j], tmp[i]
	})

	g.TeamA.Players = []*Player{}
	g.TeamB.Players = []*Player{}

	i := 0
	for _, player := range tmp {
		if i%2 == 0 {
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
	g.flow.Lock()
	defer g.flow.Unlock()

	if g.Active == false {
		log.Println("🚫 Cannot check answer for inactive games.")
		return //  cannot answer active games
	}

	t := g.GetTeam(p)
	if t == nil || t.Waiting || !t.IsPlayerAllowedToAnswer(p) {
		log.Println("wrong player answered")
		return // the wrong player attempted to answer
	}

	if !g.questionAt(t.QuestionIdx).CheckAnswer(answer.Answer) {
		g.TeamBroadcast(types.Envelope{Type: "incorrect_answer", Name: p.Name}, t)
		return
	}

	log.Println("correct answer")
	t.Turn++
	t.QuestionIdx++
	skips := t.Skips
	t.Skips = 0

	if t.Turn >= g.teamTotal(t) {
		g.EndGame(t)
		g.DisplayGameSnapshot("answered_correctly", p)
		return
	}

	g.DisplayGameSnapshot("answered_correctly", p)
	if skips == 0 {
		g.NextQuestion(t)
		return
	}
	g.startPenalty(t, time.Duration(skips)*SkipPenalty)
}

// startPenalty holds the team's next question back for d. Callers hold g.flow.
func (g *Game) startPenalty(t *Team, d time.Duration) {
	t.Waiting = true
	gen := g.generation
	payload, _ := json.Marshal(map[string]interface{}{
		"team":    t.Name,
		"seconds": int(d / time.Second),
	})
	g.Broadcast(types.Envelope{Type: "penalty_wait", Payload: payload})

	t.cancelPenalty = g.afterFunc(d, func() {
		g.flow.Lock()
		defer g.flow.Unlock()

		if gen != g.generation || !g.Active || !t.Waiting {
			return
		}
		t.Waiting = false
		t.cancelPenalty = nil
		g.NextQuestion(t)
	})
}

// handleSkip swaps the current player's question for the team's next one and queues a penalty.
func (g *Game) handleSkip(p *Player) {
	g.flow.Lock()
	defer g.flow.Unlock()

	if !g.Active {
		return
	}
	t := g.GetTeam(p)
	if t == nil || t.Waiting || !t.IsPlayerAllowedToAnswer(p) {
		return
	}

	t.Skips++
	t.QuestionIdx++

	payload, _ := json.Marshal(map[string]interface{}{
		"penalty_seconds":       int(SkipPenalty / time.Second),
		"total_penalty_seconds": int(time.Duration(t.Skips) * SkipPenalty / time.Second),
	})
	g.Broadcast(types.Envelope{Type: "question_skipped", Name: p.Name, Payload: payload})
	g.NextQuestion(t)
}

func (g *Game) handleAssignPlayerName(conn *websocket.Conn, p *Player, addPlayerPayload *AddPlayerPayload) *Player {
	if p == nil {
		p = g.AddPlayer(conn, addPlayerPayload.Name)

		payload := PlayerJoinedPayload{
			PlayerID: p.ID,
			Name:     p.Name,
		}
		g.PlayerBroadcast(types.Envelope{
			Type:    "game_player_initialized",
			GameID:  g.ID,
			Payload: utils.MustMarshal(payload),
		}, p)
	} else {
		p.Name = addPlayerPayload.Name
	}

	g.Broadcast(types.Envelope{Type: "player_joined", Name: p.Name})
	g.PlayerBroadcast(types.Envelope{Type: "joined_success", Name: p.Name}, p)

	return p
}

// WebSocket Lifecycle
func (g *Game) HandleConnection(conn *websocket.Conn, player *Player) {

	// player has entered game loop
	// boardcast they joined then
	// start reading messages

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var envelope types.Envelope
		utils.MustUnmarshal(conn, msg, &envelope)

		utils.LogPrettyJSON("🟢✉️↙️ Handle Envelope: ", envelope)

		player = g.HandleMessage(conn, player, envelope)
	}
}

func (g *Game) HandleMessage(conn *websocket.Conn, p *Player, msg types.Envelope) *Player {
	switch msg.Type {
	case "add_player":
		var addPlayerPayload AddPlayerPayload
		utils.MustUnmarshal(conn, msg.Payload, &addPlayerPayload)
		p = g.handleAssignPlayerName(conn, p, &addPlayerPayload)
	case "assign_teams":
		if p == nil {
			log.Println("🚫 Cannot assign teams before player joins.")
			return p
		}
		g.handleReassignTeams()
		g.DisplayTeamSnapshots()
		time.Sleep(50 * time.Millisecond)
		g.DisplayGameSnapshot("teams_assigned", p)
	case "reassign_teams":
		if p == nil {
			log.Println("🚫 Cannot reassign teams before player joins.")
			return p
		}
		g.handleReassignTeams()
		g.DisplayTeamSnapshots()
		time.Sleep(50 * time.Millisecond)
		g.DisplayGameSnapshot("teams_reassigned", p)
	case "show_players":
		if p == nil {
			log.Println("🚫 Cannot show players before player joins.")
			return p
		}
		g.DisplayGameSnapshot("show_players", p)
	case "check_answer":
		if p == nil {
			log.Println("🚫 Cannot answer before player joins.")
			return p
		}
		fmt.Println("msg.Payload---: ", msg.Payload)
		var answerPayload AnswerPayload
		fmt.Println("answerPayload---: ", answerPayload)
		utils.MustUnmarshal(conn, msg.Payload, &answerPayload)
		g.handleCheckAnswer(p, &answerPayload)
	case "skip_question":
		if p == nil {
			log.Println("🚫 Cannot skip before player joins.")
			return p
		}
		g.handleSkip(p)
	case "set_questions_per_player":
		if p == nil {
			log.Println("🚫 Cannot configure game before player joins.")
			return p
		}
		var payload QuestionsPerPlayerPayload
		utils.MustUnmarshal(conn, msg.Payload, &payload)
		if g.SetQuestionsPerPlayer(payload.Count) {
			g.Broadcast(types.Envelope{
				Type:    "questions_per_player_updated",
				Payload: utils.MustMarshal(QuestionsPerPlayerPayload{Count: g.QuestionsPerPlayer}),
			})
		}
	case "start":
		if p == nil {
			log.Println("🚫 Cannot start game before player joins.")
			return p
		}
		g.StartGame(p)
	case "restart_game":
		if p == nil {
			log.Println("🚫 Cannot restart game before player joins.")
			return p
		}
		g.RestartGame()
		g.DisplayGameSnapshot("game_restarted", p)
	case "update_quiz":
		if p == nil {
			log.Println("🚫 Cannot update quiz before player joins.")
			return p
		}
		var updateQuizPayload UpdateQuiz
		utils.MustUnmarshal(conn, msg.Payload, &updateQuizPayload)
		g.UpdateQuiz(p, &updateQuizPayload)
		g.DisplayGameSnapshot("quiz_updated", p)
	}

	return p
}
