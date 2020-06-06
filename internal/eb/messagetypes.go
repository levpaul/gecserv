package eb

import (
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/fb"
)

type PlayerSyncMessage struct {
	ToPlayerSID float64
	Players     []*fb.PlayerT
	Logins      []float64
	Logouts     []float64
	Tick        core.GameTick
}
