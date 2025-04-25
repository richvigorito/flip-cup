//game/team.go
package game 

type teamSnapshot struct {
    Name    string   `json:"name"`
    Turn    int      `json:"turn"`
    Players []string `json:"players"`
}

type Team struct {
	Players	[]*Player
	Name		string
	Turn		int
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

func (t *Team) snapshot() teamSnapshot {
    return teamSnapshot{
        Name:    t.Name,
        Turn:    t.Turn,
        Players: t.ExtractPlayerNames(),
    }
}
