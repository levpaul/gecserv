package main

import (
	"encoding/json"
	"errors"
	"github.com/levpaul/idolscape-backend/internal/cmdflags"
	"github.com/levpaul/idolscape-backend/internal/debug"
	"github.com/levpaul/idolscape-backend/internal/flusher"
	"github.com/levpaul/idolscape-backend/internal/ingest"
	"github.com/levpaul/idolscape-backend/internal/network"
	"github.com/levpaul/idolscape-backend/internal/propagation"
	"github.com/levpaul/idolscape-backend/internal/validation"
	"github.com/levpaul/scratch/pkg/signal"
	"github.com/pion/webrtc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	cmdflags.Parse()
}

var pipelineErrCh = make(chan error)

func main() {
	if *cmdflags.DevMode == true {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		go debug.StartDebugServer()
	}

	go network.StartNetworkManager()
	go startWebServer()

	startPipeline("ingest", ingest.Start)
	startPipeline("validation", validation.Start)
	startPipeline("flusher", flusher.Start)
	startPipeline("propagation", propagation.Start)

	select {
	case err := <-pipelineErrCh:
		log.Err(err).Send()
		return
	}
}

func startPipeline(plName string, pipeline func(<-chan error)) {
	go func() {
		pipeline(pipelineErrCh)
		pipelineErrCh <- errors.New("pipeline job returned unexpectedly")
	}()
}

func startWebServer() {
	server := http.NewServeMux()
	server.HandleFunc("/new_rtc_session", newRTCSessionHandler)
	addr := "0.0.0.0:8899"
	log.Info().Msg("Start web server on " + addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		panic(err.Error())
	}
}

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
	var conn *network.Connection
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
		conn = network.NewConnection(peerConnection, dataChannel)
		conn.SendInitState()
	})

	dataChannel.OnClose(func() {
		conn.Disconnect()
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {

		messageType := struct{ Type string }{}
		err := json.Unmarshal(msg.Data, &messageType)
		if err != nil {
			log.Err(err).Str("Message Data", string(msg.Data)).Msg("Error unmarshalling message from client")
			return
		}

		if messageType.Type == "getchars" {
			conn.SendOtherCharsState()
		}
	})
	return dataChannel, nil
}
