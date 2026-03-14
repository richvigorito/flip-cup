package game

import (
	"testing"
	"flip-cup/internal/quiz"
)

// newTestQuestionFile returns a minimal QuestionFile suitable for game unit tests.
func newTestQuestionFile() *quiz.QuestionFile {
	return &quiz.QuestionFile{
		Filename: "test.yaml",
		Questions: []*quiz.Question{
			{Prompt: "Q1", Answers: []string{"A1"}},
			{Prompt: "Q2", Answers: []string{"A2"}},
		},
	}
}

// --- Game.AddPlayer ---
// AddPlayer distributes players evenly across TeamA and TeamB (balancing by count).

func TestGame_AddPlayer_BalancesTeams(t *testing.T) {
	g := NewGame(newTestQuestionFile())

	// First player goes to TeamA (both teams empty, TeamB.len == TeamA.len == 0)
	p1 := g.AddPlayer(nil, "Alice")
	if g.GetTeam(p1) == nil {
		t.Fatal("Player 1 should belong to a team")
	}

	// Second player should go to the other team
	p2 := g.AddPlayer(nil, "Bob")
	team1 := g.GetTeam(p1)
	team2 := g.GetTeam(p2)
	if team1 == team2 {
		t.Error("Expected players to be on different teams with two players total")
	}

	// Third player: teams have 1 player each, so goes to TeamA
	p3 := g.AddPlayer(nil, "Carol")
	team3 := g.GetTeam(p3)
	if team3 == nil {
		t.Error("Third player should be assigned to a team")
	}

	// Total of 3 players should be distributed across both teams
	total := len(g.TeamA.Players) + len(g.TeamB.Players)
	if total != 3 {
		t.Errorf("Expected 3 total players, got %d", total)
	}
}

// --- Game.RemovePlayer ---

func TestGame_RemovePlayer(t *testing.T) {
	g := NewGame(newTestQuestionFile())

	p1 := g.AddPlayer(nil, "Alice")
	p2 := g.AddPlayer(nil, "Bob")

	g.RemovePlayer(p1)

	// p1 should be gone from its team
	if g.GetTeam(p1) != nil {
		t.Error("Removed player should not belong to any team")
	}

	// p2 should still be present
	if g.GetTeam(p2) == nil {
		t.Error("Non-removed player should still belong to a team")
	}

	total := len(g.TeamA.Players) + len(g.TeamB.Players)
	if total != 1 {
		t.Errorf("Expected 1 player remaining, got %d", total)
	}
}

// --- Game.GetTeam ---

func TestGame_GetTeam_ReturnsCorrectTeam(t *testing.T) {
	g := NewGame(newTestQuestionFile())

	p1 := g.AddPlayer(nil, "Alice")
	p2 := g.AddPlayer(nil, "Bob")

	team1 := g.GetTeam(p1)
	team2 := g.GetTeam(p2)

	if team1 == nil {
		t.Error("Expected team for p1")
	}
	if team2 == nil {
		t.Error("Expected team for p2")
	}
	// The two players should end up on different teams
	if team1 == team2 {
		t.Error("Expected p1 and p2 to be on different teams")
	}
}

func TestGame_GetTeam_ReturnsNilForUnknownPlayer(t *testing.T) {
	g := NewGame(newTestQuestionFile())
	ghost := newTestPlayer("unknown", "Ghost")

	if g.GetTeam(ghost) != nil {
		t.Error("Expected nil for player not in any team")
	}
}

// --- Game.ReconnectPlayer ---

func TestGame_ReconnectPlayer_Found(t *testing.T) {
	g := NewGame(newTestQuestionFile())
	p := g.AddPlayer(nil, "Alice")

	reconnected, team := g.ReconnectPlayer(p.ID, nil)

	if reconnected == nil {
		t.Fatal("Expected to find the player for reconnection")
	}
	if reconnected.ID != p.ID {
		t.Errorf("Expected player ID %s, got %s", p.ID, reconnected.ID)
	}
	if team == nil {
		t.Error("Expected team to be returned on reconnect")
	}
}

func TestGame_ReconnectPlayer_NotFound(t *testing.T) {
	g := NewGame(newTestQuestionFile())
	g.AddPlayer(nil, "Alice")

	reconnected, team := g.ReconnectPlayer("nonexistent-id", nil)

	if reconnected != nil {
		t.Error("Expected nil player for unknown ID")
	}
	if team != nil {
		t.Error("Expected nil team for unknown ID")
	}
}

// --- Game.RestartGame ---

func TestGame_RestartGame_ResetsState(t *testing.T) {
	g := NewGame(newTestQuestionFile())
	g.AddPlayer(nil, "Alice")
	g.AddPlayer(nil, "Bob")

	// Simulate an in-progress game
	g.Active = true
	g.TeamA.Turn = 2
	g.TeamB.Turn = 3

	g.RestartGame()

	if g.Active {
		t.Error("Expected game to be inactive after restart")
	}
	if g.TeamA.Turn != 0 {
		t.Errorf("Expected TeamA.Turn to be 0 after restart, got %d", g.TeamA.Turn)
	}
	if g.TeamB.Turn != 0 {
		t.Errorf("Expected TeamB.Turn to be 0 after restart, got %d", g.TeamB.Turn)
	}
}

// --- Game.handleReassignTeams ---

func TestGame_HandleReassignTeams_RedistributesAllPlayers(t *testing.T) {
	g := NewGame(newTestQuestionFile())

	// Add 4 players (2 per team initially)
	g.AddPlayer(nil, "Alice")
	g.AddPlayer(nil, "Bob")
	g.AddPlayer(nil, "Carol")
	g.AddPlayer(nil, "Dave")

	g.handleReassignTeams()

	total := len(g.TeamA.Players) + len(g.TeamB.Players)
	if total != 4 {
		t.Errorf("Expected 4 players after reassign, got %d", total)
	}
}

func TestGame_HandleReassignTeams_BlockedWhenActive(t *testing.T) {
	g := NewGame(newTestQuestionFile())

	p1 := g.AddPlayer(nil, "Alice")
	g.AddPlayer(nil, "Bob")

	// Force a known distribution before activating
	teamBefore := g.GetTeam(p1)
	g.Active = true

	g.handleReassignTeams()

	// Teams should not have changed
	teamAfter := g.GetTeam(p1)
	if teamBefore != teamAfter {
		t.Error("Expected teams to remain unchanged when game is active")
	}
}
