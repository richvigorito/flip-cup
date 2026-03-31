package game

import (
	"testing"
	"time"

	"flip-cup/internal/quiz"
)

func TestCleanupStaleGames(t *testing.T) {
	gm := NewGameManager()
	qf := &quiz.QuestionFile{
		Filename:  "test.yaml",
		Questions: []*quiz.Question{},
	}

	fresh := gm.NewGame(qf)
	stale := gm.NewGame(qf)

	stale.mu.Lock()
	stale.LastActivity = time.Now().Add(-2 * time.Hour)
	stale.mu.Unlock()

	gm.CleanupStaleGames(time.Hour)

	if gm.GetGame(fresh.ID) == nil {
		t.Fatalf("fresh game %s should still exist", fresh.ID)
	}

	if gm.GetGame(stale.ID) != nil {
		t.Fatalf("stale game %s should have been deleted", stale.ID)
	}
}

func TestPruneStaleGamesReturnsDeletedIDs(t *testing.T) {
	gm := NewGameManager()
	qf := &quiz.QuestionFile{Filename: "test.yaml"}

	fresh := gm.NewGame(qf)
	stale := gm.NewGame(qf)
	stale.LastActivity = time.Now().Add(-(30*time.Minute + 5*time.Minute))

	deleted := gm.PruneStaleGames(30 * time.Minute)

	if len(deleted) != 1 {
		t.Fatalf("expected 1 deleted game, got %d", len(deleted))
	}

	if deleted[0] != stale.ID {
		t.Fatalf("expected stale game ID %s, got %s", stale.ID, deleted[0])
	}

	if gm.GetGame(fresh.ID) == nil {
		t.Fatalf("fresh game %s should still exist", fresh.ID)
	}
}

func TestGetStaleGames(t *testing.T) {
	gm := NewGameManager()
	qf := &quiz.QuestionFile{Filename: "test.yaml"}

	fresh := gm.NewGame(qf)
	stale := gm.NewGame(qf)

	stale.mu.Lock()
	stale.LastActivity = time.Now().Add(-40 * time.Minute)
	stale.mu.Unlock()

	staleGames := gm.GetStaleGames(30 * time.Minute)

	if len(staleGames) != 1 {
		t.Fatalf("expected 1 stale game, got %d", len(staleGames))
	}

	if staleGames[0].ID != stale.ID {
		t.Fatalf("expected stale game ID %s, got %s", stale.ID, staleGames[0].ID)
	}

	for _, candidate := range staleGames {
		if candidate.ID == fresh.ID {
			t.Fatalf("fresh game %s should not be marked stale", fresh.ID)
		}
	}
}
