package netconn

import (
	"fmt"
	"github.com/levpaul/idolscape-backend/internal/eb"
)

var pipeErr chan<- error

func Start(pErr chan<- error) error {
	pipeErr = pErr
	go start()
	return nil
}

func start() {
	nc := make(chan eb.Event, 20)
	eb.Subscribe(eb.NCONNECT, nc)

	for {
		select {
		case conn := <-nc:
			x := conn.Data.(eb.NetworkConnection)
			fmt.Printf("Got network conn: %#v\n", x)
		}
	}
}

//
//func NewConnection(pc *webrtc.PeerConnection, dc *webrtc.DataChannel) {
//	event.Publish(event.NCONNECT, event.NetworkConnection{pc, dc})
//}
//connsLock.Lock()
//defer connsLock.Unlock()
//conns = append(conns, Connection{
//Uuid: uuid2.New(),
//pc:   pc,
//dc:   dc,
//GS:   state.NewGameState(),
//})
//
//newConnLock.Lock()
//defer newConnLock.Unlock()
//newConns = append(newConns, &conns[len(conns)-1])
//
//return &conns[len(conns)-1]
