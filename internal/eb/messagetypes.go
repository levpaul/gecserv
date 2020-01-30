package eb

import (
	uuid2 "github.com/google/uuid"
	"github.com/pion/webrtc"
)

type NetworkConnection struct {
	AID uuid2.UUID // Avatar ID
	DC  *webrtc.DataChannel
	// Add login token here
}
