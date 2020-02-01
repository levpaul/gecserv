package ingest

import (
	"encoding/json"
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/levpaul/idolscape-backend/pkg/signal"
	"github.com/pion/webrtc"
	"github.com/rs/zerolog/log"
	"io/ioutil"
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

	dataChannel, err := initPeerConnection(peerConnection)
	if err != nil {
		panic(err)
	}
	log.Info().Interface("DataChannel", dataChannel).Send()

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
	var aid uuid2.UUID
	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		return nil, err
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			log.Info().Msg("Closing connection")
			dataChannel.Close()
			peerConnection.Close()
		}
	})

	dataChannel.OnOpen(func() {
		aid = uuid2.New() // TODO: Replace with login / persistance
		eb.Publish(eb.Event{eb.N_CONNECT, eb.N_CONNECT_T(&core.AvatarPubConn{
			AID:  aid,
			Conn: dataChannel,
		})})
	})

	dataChannel.OnClose(func() {
		log.Info().Msg("Disconnecting dc")
		eb.Publish(eb.Event{
			Topic: eb.N_DISCONN,
			Data:  eb.N_DISCONN_T(&aid),
		})
	})

	// Register text message handling -TODO: Make this publish to validation topic
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		messageType := struct{ Type string }{}
		err := json.Unmarshal(msg.Data, &messageType)
		if err != nil {
			log.Err(err).Str("Message Data", string(msg.Data)).Msg("Error unmarshalling message from client")
			return
		}

		if messageType.Type == "getchars" {
			log.Info().Msg("Sending other chars state")
			//conn.SendOtherCharsState()
		}
	})
	return dataChannel, nil
}
