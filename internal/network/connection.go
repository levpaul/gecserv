package network

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/game"
	"github.com/pion/webrtc"
	"log"
	"sync"
)

var conns []Connection
var connsLock sync.Mutex

var newConns []*Connection
var newConnLock sync.Mutex

type Connection struct {
	Uuid uuid2.UUID
	pc   *webrtc.PeerConnection
	dc   *webrtc.DataChannel
	GS   *game.CharState
}

func NewConnection(pc *webrtc.PeerConnection, dc *webrtc.DataChannel) *Connection {
	connsLock.Lock()
	defer connsLock.Unlock()
	conns = append(conns, Connection{
		Uuid: uuid2.New(),
		pc:   pc,
		dc:   dc,
		GS:   game.NewGameState(),
	})

	newConnLock.Lock()
	defer newConnLock.Unlock()
	newConns = append(newConns, &conns[len(conns)-1])

	return &conns[len(conns)-1]
}

func (c *Connection) SendNewPlayer(color int32, x float32, y float32) {
	c.dc.SendText(fmt.Sprintf(`{"type": "newchar", "color": %d, "x":%f,"y":%f}`, color, x, y))
}

func (c *Connection) Disconnect() {
	log.Println("Disconnecting client - Levi - UUID:", c.Uuid)
	connsLock.Lock()
	defer connsLock.Unlock()

	for i := range conns {
		if conns[i].Uuid == c.Uuid {
			conns[i] = conns[len(conns)-1]
			conns = conns[:len(conns)-1]
			return
		}
	}
}

// ===========================================================================

type MoveMessage struct {
	X float32
	Y float32
}
