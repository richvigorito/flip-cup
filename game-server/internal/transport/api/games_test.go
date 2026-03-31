package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"flip-cup/internal/game"
	"flip-cup/internal/quiz"
)

func TestDeleteStaleGamesPrunesStaleGames(t *testing.T) {
	manager := game.NewGameManager()
	qf := &quiz.QuestionFile{Filename: "test.yaml"}

	fresh := manager.NewGame(qf)
	stale := manager.NewGame(qf)
	stale.LastActivity = time.Now().Add(-(30*time.Minute + 5*time.Minute))

	req := httptest.NewRequest(http.MethodDelete, "/games/stale", nil)
	rec := httptest.NewRecorder()

	deleteStaleGames(manager, 30*time.Minute).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var response deleteGamesResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Status != "stale" {
		t.Fatalf("expected status stale, got %s", response.Status)
	}

	if response.DeletedCount != 1 {
		t.Fatalf("expected 1 deleted game, got %d", response.DeletedCount)
	}

	if len(response.DeletedIDs) != 1 || response.DeletedIDs[0] != stale.ID {
		t.Fatalf("expected deleted stale ID %s, got %v", stale.ID, response.DeletedIDs)
	}

	if manager.GetGame(stale.ID) != nil {
		t.Fatalf("stale game %s should have been deleted", stale.ID)
	}

	if manager.GetGame(fresh.ID) == nil {
		t.Fatalf("fresh game %s should still exist", fresh.ID)
	}
}
