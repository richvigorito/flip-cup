//game/player.go
package game 

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	ID   string
	Name string
	Conn *websocket.Conn
//	Index int
}
