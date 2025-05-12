//game/player.go
package game 

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	ID   string
	Name string
	Conn *websocket.Conn
}

func NewPlayer(conn *websocket.Conn, name string) * Player {
	 return &Player{
		ID:   RandID(),
		Conn: conn,
		Name: name,
	};
}
