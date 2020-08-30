package eb

import (
	"github.com/levpaul/gecserv/internal/fb"
)

// MapUpdateMsg is a message from the server to a client updating their environment
type MapUpdateMsg struct {
	ToPlayerSID float64
	Msg         fb.MapUpdateT
}

// PlayerInputMsg is a message from a player with their input for a gametick
type PlayerInputMsg struct {
	FromPlayerSID float64
	Msg           fb.PlayerInputT
}
