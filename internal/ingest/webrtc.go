package ingest

import (
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/fb"
	"github.com/levpaul/gecserv/pkg/signal"
	"github.com/pion/webrtc"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math/rand"
	"net/http"
)

func newRTCSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare the configuration - No ICE servers for now since it's all local
	config := webrtc.Configuration{}

	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	clientSD, _ := ioutil.ReadAll(r.Body)

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	_, err = initPeerConnection(peerConnection)
	if err != nil {
		panic(err)
	}

	offer := webrtc.SessionDescription{}
	signal.Decode(string(clientSD), &offer)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Output the offer in base64 so we can paste it in browser
	encodedAnswer := signal.Encode(answer)

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write([]byte(encodedAnswer))
}

func initPeerConnection(peerConnection *webrtc.PeerConnection) (*webrtc.DataChannel, error) {
	var sid float64
	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		return nil, err
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			dataChannel.Close()
			peerConnection.Close()
		}
	})

	dataChannel.OnOpen(func() {
		sid = rand.Float64() // TODO: Replace with login/persistance
		eb.Publish(eb.Event{eb.N_CONNECT, eb.N_CONNECT_T(&core.SessionPubConn{
			SID:  sid,
			Conn: dataChannel,
		})})
	})

	dataChannel.OnClose(func() {
		eb.Publish(eb.Event{
			Topic: eb.N_DISCONN,
			Data:  eb.N_DISCONN_T(sid),
		})
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Info().Str("msg", string(msg.Data)).Msg("Received message")
		pmsg := fb.GetRootAsPlayerMessage(msg.Data, 0)
		if pmsg == nil {
			log.Warn().Msg("Invalid message received")
		}
		// TODO: Send to event bus to be handled by a system
	})
	return dataChannel, nil
}
