//game/player.go
package game 

import (
	"sync"
	"github.com/gorilla/websocket"
	"flip-cup/internal/utils"
)

type Player struct {
	ID   string
	Name string
	Conn *websocket.Conn
	mu   sync.Mutex
}

type PlayerSnapshot struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func NewPlayer(conn *websocket.Conn, name string) * Player {
	 return &Player{
		ID:   utils.RandID(),
		Conn: conn,
		Name: name,
	};
}

func (p *Player) Snapshot() PlayerSnapshot {
    return PlayerSnapshot{ 
        ID: p.ID,
        Name: p.Name,
    }
}
