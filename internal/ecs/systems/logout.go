package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/levpaul/idolscape-backend/internal/ecs"
	"github.com/levpaul/idolscape-backend/internal/ecs/entities"
	"github.com/rs/zerolog/log"
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
	log.Info().Str("SID", core.SIDStr(sid)).Msg("Player logout!")
	en, ok := ls.sidsToEnts[sid]
	if !ok {
		log.Error().Str("SID", core.SIDStr(sid)).Msg("Could not find entity for session during logout")
		return
	}
	ecs.RemoveEntityFromSector(en, ls.sectorID)
}

func (ls *LogoutSystem) Add(en core.EntityID) {
	rawE, err := ecs.GetEntityFromSector(en, ls.sectorID)
	if err != nil {
		log.Err(err).Send()
	}
	e, _ := rawE.(*entities.PlayerE)
	ls.sidsToEnts[e.Sid] = en
}

func (ls *LogoutSystem) Remove(en core.EntityID) {}
