package main

import (
	"encoding/json"
	"fmt"
	"github.com/levpaul/idolscape-backend/internal/network"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/levpaul/scratch/pkg/signal"
	"github.com/pion/webrtc"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	go network.StartNetworkManager()
	go startWebServer()
	select {}
}

func startWebServer() {
	server := http.NewServeMux()
	server.HandleFunc("/new_rtc_session", newRTCSessionHandler)
	addr := "0.0.0.0:8899"
	log.Println("Start web server on ", addr)
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
	log.Println("DataChan: ", dataChannel)

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
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			log.Println("Closing data channel, err: ", dataChannel.Close())
			log.Println("Closing peer connection, err: ", peerConnection.Close())
		}
	})

	dataChannel.OnOpen(func() {
		conn = network.NewConnection(peerConnection, dataChannel)
		log.Printf("New connection opened; uuid: '%s'\n", conn.Uuid)
		conn.SendInitState()
	})

	dataChannel.OnClose(func() {
		conn.Disconnect()
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("[%s] Message from DataChannel '%s': '%s'\n", conn.Uuid, dataChannel.Label(), string(msg.Data))

		messageType := struct{ Type string }{}
		err := json.Unmarshal(msg.Data, &messageType)
		if err != nil {
			fmt.Println("Error unmarshalling messgae type: ", err.Error())
			return
		}

		if messageType.Type == "getchars" {
			conn.SendOtherCharsState()
		}
	})
	return dataChannel, nil
}
