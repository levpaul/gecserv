package netpub

import (
	"fmt"
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
	eb.Subscribe(eb.N_LOGIN_RESPONSE, nc)
	eb.Subscribe(eb.N_LOGOUT_RESPONSE, nc)

	for {
		select {
		case conn := <-nc:
			switch data := conn.Data.(type) {
			// A new network connection has been made
			case eb.N_CONNECT_T:
				handleNewNetworkConn(data)

			// The network connection has been interrupted
			case eb.N_DISCONN_T:
				handleNetworkDisconnection(data)

			// An update from the server needs to be sent to a player
			case eb.N_PLAYER_SYNC_T:
				handlePlayerSync(data)
				break

			// Send back a successful login response to client with some bootstrapping information
			case eb.N_LOGIN_RESPONSE_T:
				handleLoginResponse(data)

			// Send back a successful logout response to client with some helper information - TODO: Currently not called anywhere
			case eb.N_LOGOUT_RESPONSE_T:
				handleLogoutResponse(data)

			default:
				log.Error().Int("type", int(conn.Topic)).Msg("Unsupported message type received at netpub")
				fmt.Printf("DATA: %v", data)
			}
		}
	}
}

func handleNewNetworkConn(data eb.N_CONNECT_T) {
	p := generateNewCharacter(data.SID) // TODO: Replace with persistence fetching
	// use aPConn.AID later on for bookkeeping connections to logins
	pConnMap[p.Sid] = playerConn{
		conn: data.Conn,
		p:    p,
	}
	// We have established WebRTC + PlayerLogin (AID) + PlayerObject now, publish to Simulator
	eb.Publish(eb.Event{
		Topic: eb.S_LOGIN,
		Data:  eb.S_LOGIN_T(p),
	})
}

func handleNetworkDisconnection(data eb.N_DISCONN_T) {
	log.Info().Float64("SID", float64(data)).Msg("Player disconnected")
	delete(pConnMap, float64(data))
	eb.Publish(eb.Event{
		Topic: eb.S_LOGOUT,
		Data:  eb.S_LOGOUT_T(data),
	})
}

func handlePlayerSync(data eb.N_PLAYER_SYNC_T) {
	fbBuilder.Reset()

	mapUpdate := data.Msg.Pack(fbBuilder)
	fb.ServerMessageStart(fbBuilder)
	fb.ServerMessageAddData(fbBuilder, mapUpdate)
	fb.ServerMessageAddDataType(fbBuilder, fb.ServerMessageUMapUpdate)
	message := fb.ServerMessageEnd(fbBuilder)

	fbBuilder.Finish(message)
	conn := pConnMap[data.ToPlayerSID].conn
	if err := conn.Send(fbBuilder.FinishedBytes()); err != nil {
		log.Error().Err(err).Msg("Failed to send sync update")
	}
}

func handleLoginResponse(loginRespT eb.N_LOGIN_RESPONSE_T) {
	fbBuilder.Reset()
	lrT := fb.LoginResponseT(loginRespT)
	resp := lrT.Pack(fbBuilder)

	fb.ServerMessageStart(fbBuilder)
	fb.ServerMessageAddData(fbBuilder, resp)
	fb.ServerMessageAddDataType(fbBuilder, fb.ServerMessageULoginResponse)
	message := fb.ServerMessageEnd(fbBuilder)

	fbBuilder.Finish(message)
	conn := pConnMap[loginRespT.Player.Sid].conn
	if err := conn.Send(fbBuilder.FinishedBytes()); err != nil {
		log.Error().Err(err).Msg("Failed to send login response")
	}
}

func handleLogoutResponse(logoutRespT eb.N_LOGOUT_RESPONSE_T) {
	fbBuilder.Reset()
	lrT := fb.LogoutResponseT(logoutRespT)
	resp := lrT.Pack(fbBuilder)

	fb.ServerMessageStart(fbBuilder)
	fb.ServerMessageAddData(fbBuilder, resp)
	fb.ServerMessageAddDataType(fbBuilder, fb.ServerMessageULoginResponse)
	message := fb.ServerMessageEnd(fbBuilder)

	fbBuilder.Finish(message)
	conn := pConnMap[logoutRespT.Sid].conn
	if err := conn.Send(fbBuilder.FinishedBytes()); err != nil {
		log.Error().Err(err).Msg("Failed to send logout response")
	}
}

func generateNewCharacter(sid float64) *fb.PlayerT {
	p := new(fb.PlayerT)
	p.Col = fb.Color(rand.Intn(int(fb.ColorMAXCOLOR)))
	p.Posx = (rand.Float32()) * 25
	p.Posy = (rand.Float32()) * 25
	p.Sid = sid

	return p
}
