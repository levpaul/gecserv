package eb

import "github.com/pion/webrtc"

type NetworkConnection struct {
	PC *webrtc.PeerConnection
	DC *webrtc.DataChannel
	// Add login token here
}
