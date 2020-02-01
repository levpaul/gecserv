package netpub

import (
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/levpaul/idolscape-backend/internal/fb"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
	"sync/atomic"
)

var (
	pipeErr chan<- error

	pConnMap       map[uuid2.UUID]playerConn
	sessionCounter uint64
)

// TODO: Make this struct more generic - Sender/Reciever interface
type playerConn struct {
	conn core.SenderCloser
	p    *fb.PlayerT
}

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	go startListening()
	return nil
}

func initialize() {
	pConnMap = make(map[uuid2.UUID]playerConn)
	sessionCounter = rand.Uint64()
}

func startListening() {
	nc := make(chan eb.Event, 20)
	eb.Subscribe(eb.N_CONNECT, nc)
	eb.Subscribe(eb.N_DISCONN, nc)

	for {
		select {
		case conn := <-nc:
			switch conn.Topic {

			case eb.N_CONNECT:
				aPConn := conn.Data.(eb.N_CONNECT_T)
				p := generateNewCharacter() // TODO: Replace with persistence fetching
				pConnMap[aPConn.AID] = playerConn{
					conn: aPConn.Conn,
					p:    p,
				}
				// We have established WebRTC + PlayerLogin (AID) + PlayerObject now, publish to Simulator
				eb.Publish(eb.Event{
					Topic: eb.S_LOGIN,
					Data:  eb.S_LOGIN_T(p),
				})

			case eb.N_DISCONN:
				aid := conn.Data.(eb.N_DISCONN_T)
				pid := pConnMap[*aid].p.Sid
				delete(pConnMap, *aid)
				eb.Publish(eb.Event{
					Topic: eb.S_LOGOUT,
					Data:  eb.S_LOGOUT_T(pid),
				})

			default:
				log.Error().Msg("Unsupported message type ")
			}
		}
	}
}

func generateNewCharacter() *fb.PlayerT {
	p := new(fb.PlayerT)
	p.Col = fb.ColorBlue
	p.Pos = &fb.Vec2T{
		X: (rand.Float32() - 0.5) * 40,
		Y: (rand.Float32() - 0.5) * 40,
	}
	p.Sid = math.Float64frombits(atomic.AddUint64(&sessionCounter, 1))

	return p
}