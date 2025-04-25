// player.go 
package game

type Player struct {
	ID   string
	Name string
	Conn *websocket.Conn
	Index int
}


