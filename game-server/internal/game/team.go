// internal/game/team.go
package game 

import (
	"math/rand"
	"time"
)

type TeamSnapshot struct {
    Name    string   `json:"name"`
    Turn    int      `json:"turn"`
    Players []PlayerSnapshot `json:"players"`

}

type Team struct {
	Players	[]*Player
	Name		string
	Turn		int
}

func (t *Team) AddPlayer (p *Player) {
    t.Players = append(t.Players, p)
}

func (t *Team) RemovePlayer(p *Player) {
    for i, pl := range t.Players {
		if pl == p {
			t.Players = append(t.Players[:i], t.Players[i+1:]...)
			break
		}
	}
}

func (t *Team) GetCurrentPlayer() *Player {
	return t.Players[t.Turn]
}

func (t *Team) IsPlayerAllowedToAnswer(p *Player) bool {
	return t.GetPlayerIndex(p) == t.Turn 
}

func (t *Team) GetPlayerIndex(p *Player) int {
    for i := 0; i < len(t.Players); i++ { // Fix syntax for the loop
        if p.ID == t.Players[i].ID {
            return i
        }
    }
    return -1 // Return -1 if player is not found
}

func (t *Team) GetPlayer (ID string) *Player {
	for _, player := range t.Players {
		if ID == player.ID {
			return player	
		}
	}
	return nil
}

func (t *Team) ExtractPlayerNames() []string {
    names := []string{}
    for _, p := range t.Players {
        names = append(names, p.Name)
    }
    return names
}


func (t *Team) Shuffle() {
    rand.Seed(time.Now().UnixNano()) // Seed only once
    rand.Shuffle(len(t.Players), func(i, j int) {
        t.Players[i], t.Players[j] = t.Players[j], t.Players[i]
    })
}

func (t *Team) Snapshot() TeamSnapshot {
    return TeamSnapshot{
        Name:    t.Name,
        Turn:    t.Turn,
        Players: t.ExtractPlayerSnapshots(),
    }
}

func (t *Team) ExtractPlayerSnapshots() []PlayerSnapshot {
    snapshots := make([]PlayerSnapshot, 0, len(t.Players))
    for _, p := range t.Players {
        snapshots = append(snapshots, PlayerSnapshot{
            ID:   p.ID,
            Name: p.Name,
        })
    }
    return snapshots
}
