package game

import (
	"fmt"
	"log"
	"sync"
	"time"

	"flip-cup/internal/quiz"
	"flip-cup/internal/utils"
)

type GameManager struct {
	games map[string]*Game
	mu    sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		games: make(map[string]*Game),
	}
}

func (gm *GameManager) NewGame(questionFile *quiz.QuestionFile) *Game {
	g := NewGame(questionFile)
	g.ID = utils.RandID()

	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.games[g.ID] = g

	return g
}

func (gm *GameManager) GetGame(id string) *Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	return gm.games[id]
}

func (gm *GameManager) GetAllGames() []*Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	all := make([]*Game, 0, len(gm.games))
	for _, g := range gm.games {
		all = append(all, g)
	}

	return all
}

func (gm *GameManager) DeleteGameByID(id string) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	if _, exists := gm.games[id]; !exists {
		return fmt.Errorf("game with ID %s not found", id)
	}

	delete(gm.games, id)
	return nil
}

func (gm *GameManager) GetStaleGames(maxAge time.Duration) []*Game {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	stale := []*Game{}
	for _, g := range gm.games {
		if g.IsStale(maxAge) {
			stale = append(stale, g)
		}
	}

	return stale
}

func (gm *GameManager) PruneStaleGames(maxAge time.Duration) []string {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	deleted := []string{}
	for id, g := range gm.games {
		if g.IsStale(maxAge) {
			delete(gm.games, id)
			deleted = append(deleted, id)
		}
	}

	return deleted
}

func (gm *GameManager) CleanupStaleGames(maxAge time.Duration) {
	deleted := gm.PruneStaleGames(maxAge)
	if len(deleted) == 0 {
		return
	}

	for _, id := range deleted {
		log.Printf("🧹 cleaned up stale game: %s", id)
	}
	log.Printf("🧹 cleaned up %d stale games", len(deleted))
}
