package netpub

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/fb"
	"github.com/rs/zerolog/log"
	"math/rand"
)

var (
	pipeErr chan<- error

	pConnMap  map[float64]playerConn
	fbBuilder *flatbuffers.Builder
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
	fbBuilder = flatbuffers.NewBuilder(0)
}

func startListening() {
	nc := make(chan eb.Event, 20)
	eb.Subscribe(eb.N_CONNECT, nc)
	eb.Subscribe(eb.N_DISCONN, nc)
	eb.Subscribe(eb.N_PLAYER_SYNC, nc)

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

func handlePlayerSync(syncData eb.N_PLAYER_SYNC_T) {
	fbBuilder.Reset()

	packedPlayers := []flatbuffers.UOffsetT{}
	for _, p := range syncData.Players {
		packedPlayers = append(packedPlayers, p.Pack(fbBuilder))
	}

	fb.MapUpdateStartLoginsVector(fbBuilder, len(syncData.Players))
	for i := len(packedPlayers) - 1; i >= 0; i-- {
		fbBuilder.PrependUOffsetT(packedPlayers[i])
	}
	players := fbBuilder.EndVector(len(syncData.Players))

	fb.MapUpdateStartLoginsVector(fbBuilder, len(syncData.Logins))
	for i := len(syncData.Logins) - 1; i >= 0; i-- {
		fbBuilder.PrependFloat64(syncData.Logins[i])
	}
	logins := fbBuilder.EndVector(len(syncData.Logins))

	fb.MapUpdateStartLogoutsVector(fbBuilder, len(syncData.Logouts))
	for i := len(syncData.Logouts) - 1; i >= 0; i-- {
		fbBuilder.PrependFloat64(syncData.Logouts[i])
	}
	logouts := fbBuilder.EndVector(len(syncData.Logouts))

	fb.MapUpdateStart(fbBuilder)
	fb.MapUpdateAddSeq(fbBuilder, uint32(syncData.Tick))
	fb.MapUpdateAddPsyncs(fbBuilder, players)
	fb.MapUpdateAddLogins(fbBuilder, logins)
	fb.MapUpdateAddLogouts(fbBuilder, logouts)
	mapUpdate := fb.MapUpdateEnd(fbBuilder)

	fb.ServerMessageStart(fbBuilder)
	fb.ServerMessageAddData(fbBuilder, mapUpdate)
	fb.ServerMessageAddDataType(fbBuilder, fb.ServerMessageUMapUpdate)
	message := fb.ServerMessageEnd(fbBuilder)

	fbBuilder.Finish(message)
	conn := pConnMap[syncData.ToPlayerSID].conn
	conn.Send(fbBuilder.FinishedBytes())

	//psyncload := fb.GetRootAsServerMessage(fbBuilder.FinishedBytes(), 0)
	//mpupdate := &fb.ServerMessageT{}
	//psyncload.UnPackTo(mpupdate)
	//x := mpupdate.Data.Value.(*fb.MapUpdateT)
	//
	//log.Info().Msgf("Reloaded mapupdate: %+v", x)
	//log.Info().Msgf("raw mapupdate: %+v, players: %+v", x,
	//	x.Psyncs[0])
}

func generateNewCharacter(sid float64) *fb.PlayerT {
	p := new(fb.PlayerT)
	p.Col = fb.ColorBlue
	p.Posx = (rand.Float32()) * 1000
	p.Posy = (rand.Float32()) * 1000
	p.Sid = sid

	return p
}
