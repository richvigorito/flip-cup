package game

import (
    "sync"
    "flip-cup/internal/quiz" 
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
    g.ID = RandID()

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
