package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
)

// LogoutSystem handles when a player has logged off a character for a
// given sector
type LogoutSystem struct {
	BaseSystem
	logouts    chan eb.Event
	sidsToEnts map[float64]core.EntityID
}

// TODO: The topic for logouts may need to be split per sector
func (ls *LogoutSystem) Init() {
	ls.logouts = make(chan eb.Event, 128)
	ls.sidsToEnts = make(map[float64]core.EntityID)
	eb.Subscribe(eb.S_LOGOUT, ls.logouts)
}

func (ls *LogoutSystem) Update(ctx context.Context, dt core.GameTick) {
	for {
		select {

		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func (ls *LogoutSystem) Remove(en core.EntityID) {}
