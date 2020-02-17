package netpub

import (
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/fb"
	"github.com/rs/zerolog/log"
	"math/rand"
)

var (
	pipeErr chan<- error

	pConnMap map[float64]playerConn
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
	pConnMap = make(map[float64]playerConn)
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
				p := generateNewCharacter(aPConn.SID) // TODO: Replace with persistence fetching
				// use aPConn.AID later on for bookkeeping connections to logins
				pConnMap[p.Sid] = playerConn{
					conn: aPConn.Conn,
					p:    p,
				}
				// We have established WebRTC + PlayerLogin (AID) + PlayerObject now, publish to Simulator
				eb.Publish(eb.Event{
					Topic: eb.S_LOGIN,
					Data:  eb.S_LOGIN_T(p),
				})

			case eb.N_DISCONN:
				sid := float64(conn.Data.(eb.N_DISCONN_T))
				delete(pConnMap, sid)
				eb.Publish(eb.Event{
					Topic: eb.S_LOGOUT,
					Data:  eb.S_LOGOUT_T(sid),
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

func generateNewCharacter(sid float64) *fb.PlayerT {
	p := new(fb.PlayerT)
	p.Col = fb.ColorBlue
	p.Posx = (rand.Float32()) * 1000
	p.Posy = (rand.Float32()) * 1000
	p.Sid = sid

	return p
}
