package game

import (
    "sync"
    "flip-cup/internal/quiz" 
	"flip-cup/internal/utils"
	"fmt"
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
    all := []*Game{}
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

