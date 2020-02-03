package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/levpaul/idolscape-backend/internal/ecs/entities"
	"github.com/levpaul/idolscape-backend/internal/fb"
	"github.com/rs/zerolog/log"
)

// LoginSystem handles when a player has logged into a character for a
// given sector or when a player has logged out or disconnected. It is also
// solely responsible for keeping the PlayerList for a sector up to date
type LoginSystem struct {
	BaseSystem
	loginEvents chan eb.Event
	sidsToEnts  map[float64]*entities.PlayerE
}

// TODO: The topic for loginEvents may need to be split per sector
func (ls *LoginSystem) Init() {
	ls.loginEvents = make(chan eb.Event, 128)
	eb.Subscribe(eb.S_LOGIN, ls.loginEvents)
	eb.Subscribe(eb.S_LOGOUT, ls.loginEvents)

	// Set up for player list singleton management
	ls.sidsToEnts = make(map[float64]*entities.PlayerE)
	ls.sa.SetPlayerList(ls.sidsToEnts)
}

func (ls *LoginSystem) Update(ctx context.Context, dt core.GameTick) {
	for {
		select {
		case l := <-ls.loginEvents:
			switch l.Data.(type) {
			case eb.S_LOGIN_T:
				player, ok := l.Data.(eb.S_LOGIN_T)
				if !ok {
					log.Error().Interface("data", l.Data).Msg("Failed to type assert S_LOGIN message")
					continue
				}
				ls.handleLogin(ctx, player)

			case eb.S_LOGOUT_T:
				sid, ok := l.Data.(eb.S_LOGOUT_T)
				if !ok {
					log.Error().Interface("data", l.Data).Msg("Failed to type assert S_LOGOUT message")
					continue
				}
				ls.handleLogout(ctx, float64(sid))
			default:
				log.Error().Interface("type", l.Data).Msg("Unsupported message type received on login channel")
			}
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func (ls *LoginSystem) handleLogin(ctx context.Context, player *fb.PlayerT) {
	log.Info().Str("SID", core.SIDStr(player.Sid)).Msg("New player login!")

	pEntity := &entities.PlayerE{
		BaseEntity: entities.NewBaseEntity(),
		PlayerT:    player,
	}

	ls.sa.AddEntity(pEntity)
	ls.sidsToEnts[pEntity.Sid] = pEntity
}

func (ls *LoginSystem) handleLogout(ctx context.Context, sid float64) {
	log.Info().Str("SID", core.SIDStr(sid)).Msg("Player logout!")
	en, ok := ls.sidsToEnts[sid]
	if !ok {
		log.Error().Str("SID", core.SIDStr(sid)).Msg("Could not find entity for session during logout")
		return
	}

	ls.sa.RemoveEntity(en.ID())
	delete(ls.sidsToEnts, en.Sid)
}
