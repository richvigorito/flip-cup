package game

import (
	"strings"
	"testing"
	"time"

	"flip-cup/internal/quiz"
	"flip-cup/internal/transport/types"
)

type sent struct {
	to  string
	msg types.Envelope
}

type fakeClock struct {
	pending []*fakeTimer
}

type fakeTimer struct {
	d         time.Duration
	f         func()
	cancelled bool
}

func (c *fakeClock) after(d time.Duration, f func()) func() {
	t := &fakeTimer{d: d, f: f}
	c.pending = append(c.pending, t)
	return func() { t.cancelled = true }
}

// fire runs the oldest live timer and returns its delay.
func (c *fakeClock) fire(t *testing.T) time.Duration {
	for len(c.pending) > 0 {
		ft := c.pending[0]
		c.pending = c.pending[1:]
		if !ft.cancelled {
			ft.f()
			return ft.d
		}
	}
	t.Fatal("no pending timer")
	return 0
}

func (c *fakeClock) live() int {
	n := 0
	for _, t := range c.pending {
		if !t.cancelled {
			n++
		}
	}
	return n
}

type harness struct {
	g     *Game
	clock *fakeClock
	log   []sent
	a, b  []*Player // team A and team B players
}

// newHarness builds a game with aPlayers on team A and bPlayers on team B (6 questions: Q1..Q6 -> A1..A6).
func newHarness(aPlayers, bPlayers int) *harness {
	qf := &quiz.QuestionFile{Filename: "t.yaml"}
	for _, n := range []string{"1", "2", "3", "4", "5", "6"} {
		qf.Questions = append(qf.Questions, &quiz.Question{Prompt: "Q" + n, Answers: []string{"A" + n}})
	}
	h := &harness{g: NewGame(qf), clock: &fakeClock{}}
	h.g.afterFunc = h.clock.after
	h.g.onSend = func(p *Player, m types.Envelope) { h.log = append(h.log, sent{p.Name, m}) }
	for i := 0; i < aPlayers; i++ {
		p := &Player{ID: "a" + string(rune('0'+i)), Name: "A" + string(rune('0'+i))}
		h.g.TeamA.Players = append(h.g.TeamA.Players, p)
		h.a = append(h.a, p)
	}
	for i := 0; i < bPlayers; i++ {
		p := &Player{ID: "b" + string(rune('0'+i)), Name: "B" + string(rune('0'+i))}
		h.g.TeamB.Players = append(h.g.TeamB.Players, p)
		h.b = append(h.b, p)
	}
	return h
}

func (h *harness) count(typ string) int {
	n := 0
	for _, s := range h.log {
		if s.msg.Type == typ {
			n++
		}
	}
	return n
}

// lastQuestion returns the most recent question prompt sent to the named player.
func (h *harness) lastQuestion(name string) string {
	for i := len(h.log) - 1; i >= 0; i-- {
		if h.log[i].to == name && h.log[i].msg.Type == "question" {
			return h.log[i].msg.Name
		}
	}
	return ""
}

func (h *harness) answer(p *Player) {
	q := h.lastQuestion(p.Name)
	h.g.handleCheckAnswer(p, &AnswerPayload{Answer: "A" + strings.TrimPrefix(q, "Q")})
}

func TestQuestionsPerPlayer_DefaultIsOneAndWinsAfterEveryoneAnswers(t *testing.T) {
	h := newHarness(2, 1)
	h.g.StartGame(h.a[0])

	h.answer(h.a[0])
	if h.count("winner") != 0 {
		t.Fatal("should not win after one answer on a 2-player team")
	}
	h.answer(h.a[1])
	if h.count("winner") == 0 {
		t.Fatal("team A should win after both players answer once")
	}
}

func TestQuestionsPerPlayer_TwoRotatesAndWinsAfterFour(t *testing.T) {
	h := newHarness(2, 1)
	h.g.SetQuestionsPerPlayer(2)
	h.g.StartGame(h.a[0])

	order := []*Player{h.a[0], h.a[1], h.a[0], h.a[1]}
	for i, p := range order {
		if h.count("winner") != 0 {
			t.Fatalf("won too early at step %d", i)
		}
		if h.g.TeamA.GetCurrentPlayer() != p {
			t.Fatalf("step %d: expected %s to be up", i, p.Name)
		}
		h.answer(p)
	}
	if h.count("winner") == 0 {
		t.Fatal("expected a winner after 4 correct answers")
	}
}

func TestQuestionsPerPlayer_ClampAndLockedOnceStarted(t *testing.T) {
	h := newHarness(1, 1)
	h.g.SetQuestionsPerPlayer(99)
	if h.g.QuestionsPerPlayer != MaxQuestionsPerPlayer {
		t.Fatalf("expected clamp to max, got %d", h.g.QuestionsPerPlayer)
	}
	h.g.SetQuestionsPerPlayer(0)
	if h.g.QuestionsPerPlayer != MinQuestionsPerPlayer {
		t.Fatalf("expected clamp to min, got %d", h.g.QuestionsPerPlayer)
	}
	h.g.SetQuestionsPerPlayer(3)
	h.g.StartGame(h.a[0])
	if h.g.SetQuestionsPerPlayer(1) || h.g.QuestionsPerPlayer != 3 {
		t.Fatal("setting must be rejected after start")
	}
}

func TestSkip_SendsDifferentQuestionAndNotifiesEveryone(t *testing.T) {
	h := newHarness(2, 1)
	h.g.StartGame(h.a[0])
	before := h.lastQuestion("A0")

	h.g.handleSkip(h.a[0])

	after := h.lastQuestion("A0")
	if after == before {
		t.Fatalf("expected a different question, got %s again", after)
	}
	if h.g.TeamA.GetCurrentPlayer() != h.a[0] {
		t.Fatal("skipping must not advance the turn")
	}
	recipients := map[string]bool{}
	for _, s := range h.log {
		if s.msg.Type == "question_skipped" {
			if s.msg.Name != "A0" {
				t.Fatalf("skip should name the skipper, got %q", s.msg.Name)
			}
			recipients[s.to] = true
		}
	}
	for _, p := range append(append([]*Player{}, h.a...), h.b...) {
		if !recipients[p.Name] {
			t.Errorf("%s was not notified of the skip", p.Name)
		}
	}
}

func TestSkip_PenaltyScalesWithSkips(t *testing.T) {
	cases := []struct {
		skips int
		want  time.Duration
	}{{0, 0}, {1, SkipPenalty}, {2, 2 * SkipPenalty}}

	for _, c := range cases {
		h := newHarness(2, 1)
		h.g.StartGame(h.a[0])
		for i := 0; i < c.skips; i++ {
			h.g.handleSkip(h.a[0])
		}
		h.answer(h.a[0])

		if c.skips == 0 {
			if h.clock.live() != 0 {
				t.Fatal("no penalty expected without skips")
			}
			if h.lastQuestion("A1") == "" {
				t.Fatal("next player should be asked immediately")
			}
			continue
		}
		if h.lastQuestion("A1") != "" {
			t.Fatalf("skips=%d: next player asked before the penalty elapsed", c.skips)
		}
		if got := h.clock.fire(t); got != c.want {
			t.Fatalf("skips=%d: penalty %v, want %v", c.skips, got, c.want)
		}
		if h.lastQuestion("A1") == "" {
			t.Fatalf("skips=%d: next player not asked after penalty", c.skips)
		}
	}
}

func TestSkip_RejectedWhenNotAllowed(t *testing.T) {
	h := newHarness(2, 1)

	h.g.handleSkip(h.a[0]) // inactive
	if h.count("question_skipped") != 0 {
		t.Fatal("skip must be rejected before the game starts")
	}

	h.g.StartGame(h.a[0])
	h.g.handleSkip(h.a[1]) // not their turn
	if h.count("question_skipped") != 0 {
		t.Fatal("skip must be rejected for a player who is not up")
	}

	h.g.handleSkip(h.a[0])
	h.answer(h.a[0]) // now team A is waiting out a penalty
	skipsSoFar := h.count("question_skipped")
	h.g.handleSkip(h.a[1])
	if h.count("question_skipped") != skipsSoFar {
		t.Fatal("skip must be rejected while the team is waiting")
	}
}

func TestSkip_RestartCancelsPendingQuestion(t *testing.T) {
	h := newHarness(2, 1)
	h.g.StartGame(h.a[0])
	h.g.handleSkip(h.a[0])
	h.answer(h.a[0])

	h.g.RestartGame()
	if h.clock.live() != 0 {
		t.Fatal("restart should cancel the pending penalty timer")
	}
}

func TestSkip_RestartedTimerDoesNotAskStalePlayer(t *testing.T) {
	h := newHarness(2, 1)
	h.g.StartGame(h.a[0])
	h.g.handleSkip(h.a[0])
	h.answer(h.a[0])
	stale := h.clock.pending[0]

	h.g.RestartGame()
	stale.f() // even if the callback fires anyway, the generation check stops it

	if h.lastQuestion("A1") != "" {
		t.Fatal("stale penalty callback asked a question after restart")
	}
}

func TestQuestionIndex_WrapsInsteadOfPanicking(t *testing.T) {
	h := newHarness(1, 1)
	h.g.StartGame(h.a[0])
	for i := 0; i < 10; i++ { // more skips than questions
		h.g.handleSkip(h.a[0])
	}
	if h.lastQuestion("A0") == "" {
		t.Fatal("expected a question after many skips")
	}
}

func TestSkip_PenaltyDelaysWinOnFinalQuestion(t *testing.T) {
	h := newHarness(1, 1)
	h.g.StartGame(h.a[0])
	h.g.handleSkip(h.a[0])
	h.g.handleSkip(h.a[0])
	h.answer(h.a[0])

	if h.count("winner") != 0 {
		t.Fatal("win must wait out the skip penalty")
	}
	if got := h.clock.fire(t); got != 2*SkipPenalty {
		t.Fatalf("penalty %v, want %v", got, 2*SkipPenalty)
	}
	if h.count("winner") == 0 {
		t.Fatal("expected a winner after the penalty")
	}
}

func TestSkip_OtherTeamWinningCancelsPendingWin(t *testing.T) {
	h := newHarness(1, 1)
	h.g.StartGame(h.a[0])
	h.g.handleSkip(h.a[0])
	h.answer(h.a[0]) // A is waiting out a penalty before winning
	h.answer(h.b[0]) // B wins first
	winners := h.count("winner")
	h.clock.fire(t)
	if h.count("winner") != winners {
		t.Fatal("pending win must not fire after the game ended")
	}
}

func TestQuestionsPerPlayer_StaysLockedAfterGameOverUntilRestart(t *testing.T) {
	h := newHarness(1, 1)
	h.g.StartGame(h.a[0])
	h.answer(h.a[0]) // team A wins, game inactive but round not reset

	if h.g.SetQuestionsPerPlayer(3) {
		t.Fatal("setting must stay locked on the game-over screen")
	}
	h.g.RestartGame()
	if !h.g.SetQuestionsPerPlayer(3) {
		t.Fatal("setting should unlock after restart")
	}
}

func TestSkip_PenaltyEndedBroadcastToEveryone(t *testing.T) {
	h := newHarness(2, 1)
	h.g.StartGame(h.a[0])
	h.g.handleSkip(h.a[0])
	h.answer(h.a[0])
	h.clock.fire(t)
	if h.count("penalty_ended") != 3 {
		t.Fatalf("expected penalty_ended for all 3 players, got %d", h.count("penalty_ended"))
	}
}
