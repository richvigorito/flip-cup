package game

import (
	"testing"
	"time"
	"flip-cup/internal/quiz"
)

func TestCleanupStaleGames(t *testing.T) {
	// 1. Setup
	gm := NewGameManager()
	
	// Create a minimal QuestionFile structure
	qf := &quiz.QuestionFile{
		Filename: "test.yaml",
		// We don't need actual questions for this test
		Questions: []*quiz.Question{},
	}

	// 2. Create games
	g1 := gm.NewGame(qf)
	g2 := gm.NewGame(qf)

	// 3. Manipulate LastActivity
	// g1 is active now (default NewGame sets to Now())
	// g2 should be stale. We set it to 2 hours ago.
	g2.mu.Lock()
	g2.LastActivity = time.Now().Add(-2 * time.Hour)
	g2.mu.Unlock()

	// 4. Run Cleanup with 1 hour threshold
	// g2 is 2 hours old, so it > 1 hour, should be deleted.
	gm.CleanupStaleGames(1 * time.Hour)

	// 5. Verify
	if gm.GetGame(g1.ID) == nil {
		t.Errorf("Game 1 (active) should still exist")
	}
	if gm.GetGame(g2.ID) != nil {
		t.Errorf("Game 2 (stale) should have been deleted")
	}
}

func TestGetStaleGames(t *testing.T) {
	gm := NewGameManager()
	qf := &quiz.QuestionFile{Filename: "test.yaml"}

	g1 := gm.NewGame(qf)
	g2 := gm.NewGame(qf)

	// Make g2 stale
	g2.mu.Lock()
	g2.LastActivity = time.Now().Add(-40 * time.Minute)
	g2.mu.Unlock()

	// Threshold 30m
	stale := gm.GetStaleGames(30 * time.Minute)

	if len(stale) != 1 {
		t.Errorf("Expected 1 stale game, got %d", len(stale))
	} else if stale[0].ID != g2.ID {
		t.Errorf("Expected stale game ID %s, got %s", g2.ID, stale[0].ID)
	}

    // Ensure g1 (active) is not returned
    for _, s := range stale {
        if s.ID == g1.ID {
            t.Errorf("Active game %s should not be marked stale", g1.ID)
        }
    }
}
