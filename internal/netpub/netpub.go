package netpub

import (
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/fb"
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

			case eb.N_PLAYER_SYNC:
				pSyncData, ok := conn.Data.(eb.N_PLAYER_SYNC_T)
				if !ok {
					log.Error().Msgf("Non playersync data sent to eventbus on player sync channel - %#v", conn.Data)
					continue
				}
				handlePlayerSync(pSyncData)

			default:
				log.Error().Msg("Unsupported message type ")
			}
		}
	}
}

func handlePlayerSync(conn eb.N_PLAYER_SYNC_T) {

}

func generateNewCharacter() *fb.PlayerT {
	p := new(fb.PlayerT)
	p.Col = fb.ColorBlue
	p.Posx = (rand.Float32()) * 1000
	p.Posy = (rand.Float32()) * 1000
	p.Sid = math.Float64frombits(atomic.AddUint64(&sessionCounter, 1))

	return p
}
