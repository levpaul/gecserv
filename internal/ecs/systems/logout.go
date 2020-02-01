package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/rs/zerolog/log"
)

// LogoutSystem handles when a player has logged off a character for a
// given sector
type LogoutSystem struct {
	logouts chan eb.Event
}

// TODO: The topic for logouts may need to be split per sector
func (ls *LogoutSystem) Init() {
	ls.logouts = make(chan eb.Event, 128)
	eb.Subscribe(eb.S_LOGOUT, ls.logouts)
}

func (ls *LogoutSystem) Update(ctx context.Context, dt core.GameTick) {
	for {
		select {
		case l := <-ls.logouts:
			sid, ok := l.Data.(eb.S_LOGOUT_T)
			if !ok {
				log.Error().Interface("data", l.Data).Msg("Failed to type assert S_LOGOUT message")
				continue
			}
			ls.handleLogout(ctx, float64(sid))
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func (ls *LogoutSystem) handleLogout(ctx context.Context, sid float64) {
	log.Info().Float64("SID", sid).Msg("Player logout!")
}
