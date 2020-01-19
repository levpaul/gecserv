package network

import (
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/game"
	"github.com/pion/webrtc"
	"sync"
)

var conns []Connection
var newConnLock sync.Mutex

type Connection struct {
	Uuid uuid2.UUID
	pc   *webrtc.PeerConnection
	dc   *webrtc.DataChannel
	GS   *game.GameState
}

func NewConnection(pc *webrtc.PeerConnection, dc *webrtc.DataChannel) *Connection {
	newConnLock.Lock()
	defer newConnLock.Unlock()
	conns = append(conns, Connection{
		Uuid: uuid2.New(),
		pc:   pc,
		dc:   dc,
		GS:   game.NewGameState(),
	})
	return &conns[len(conns)-1]
}

// ===========================================================================

type MoveMessage struct {
	X float32
	Y float32
}
