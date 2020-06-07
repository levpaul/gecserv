package eb

import (
	"github.com/levpaul/gecserv/internal/fb"
)

type MapUpdateMsg struct {
	ToPlayerSID float64
	Msg         fb.MapUpdateT
}
