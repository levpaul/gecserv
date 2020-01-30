package netconn

import (
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/levpaul/idolscape-backend/internal/fb"
	"github.com/pion/webrtc"
	"math/rand"
)

var (
	pipeErr chan<- error

	pConnMap map[uuid2.UUID]playerConn
)

type playerConn struct {
	dc *webrtc.DataChannel
	pc *webrtc.PeerConnection
	p  *fb.PlayerT
}

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	go start()
	return nil
}
func initialize() {
	pConnMap = make(map[uuid2.UUID]playerConn)
}

func start() {
	nc := make(chan eb.Event, 20)
	eb.Subscribe(eb.N_CONNECT, nc)

	for {
		select {
		case conn := <-nc:
			tConn := conn.Data.(eb.NetworkConnection)
			p := generateNewCharacter() // TODO: Replace with persistence fetching
			pConnMap[tConn.AID] = playerConn{
				dc: tConn.DC,
				pc: tConn.PC,
				p:  p,
			}
			// We have established WebRTC + PlayerLogin (AID) + PlayerObject now, publish to Simulator
			eb.Publish(eb.Event{
				Topic: eb.S_LOGIN,
				Data:  p,
			})
		}
	}
}

func generateNewCharacter() *fb.PlayerT {
	p := new(fb.PlayerT)
	p.Col = fb.ColorBlue
	p.Pos.X = (rand.Float32() - 0.5) * 40
	p.Pos.Y = (rand.Float32() - 0.5) * 40

	return p
}

//
//func NewConnection(pc *webrtc.PeerConnection, dc *webrtc.DataChannel) {
//	event.Publish(event.N_CONNECT, event.NetworkConnection{pc, dc})
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
