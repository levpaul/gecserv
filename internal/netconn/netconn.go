package netconn

import (
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/state"
	"github.com/pion/webrtc"
)

var pipeErr chan<- error

func Start(pErr chan<- error) error {
	pipeErr = pErr
	go start()
	return nil
}

func start() {
	for {
		select {}
	}
}

type NetworkConnectEvent {

}


func NewConnection(pc *webrtc.PeerConnection, dc *webrtc.DataChannel) *Connection {
	connsLock.Lock()
	defer connsLock.Unlock()
	conns = append(conns, Connection{
		Uuid: uuid2.New(),
		pc:   pc,
		dc:   dc,
		GS:   state.NewGameState(),
	})

	newConnLock.Lock()
	defer newConnLock.Unlock()
	newConns = append(newConns, &conns[len(conns)-1])

	return &conns[len(conns)-1]
}
