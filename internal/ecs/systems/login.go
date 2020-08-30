package systems

import (
	"context"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/ecs/components"
	"github.com/levpaul/gecserv/internal/ecs/entities"
	"github.com/levpaul/gecserv/internal/fb"
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

func (ls *LoginSystem) Init() {
	ls.loginEvents = make(chan eb.Event, 128)
	eb.Subscribe(eb.S_LOGIN, ls.loginEvents)
	eb.Subscribe(eb.S_LOGOUT, ls.loginEvents)

	// Set up for player list singleton management
	ls.sidsToEnts = make(map[float64]*entities.PlayerE)
	ls.sa.SetPlayerListSingleton(ls.sidsToEnts)
}

func (ls *LoginSystem) Update(ctx context.Context, dt core.GameTick) {
	for {
		select {
		case l := <-ls.loginEvents:
			switch data := l.Data.(type) {
			// Network login event
			case eb.S_LOGIN_T:
				ls.handleLogin(ctx, data)

			// Network disconnect event
			case eb.S_LOGOUT_T:
				ls.handleDisconnect(ctx, float64(data))

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
		Position:   components.Position{player.Posx, player.Posy},
		Momentum:   components.Momentum{0, 0},
		Changeable: components.Changeable{true},
		NetworkedSession: components.NetworkedSession{
			Sid:       player.Sid,
			LoginTick: ls.sa.GetSectorTick(),
		},
		StateHistory: components.StateHistory{},
		Colored:      components.Colored{player.Col},
	}

	ls.sa.AddEntity(pEntity)
	ls.sidsToEnts[pEntity.Sid] = pEntity

	eb.Publish(eb.Event{
		Topic: eb.N_LOGIN_RESPONSE,
		Data: eb.N_LOGIN_RESPONSE_T{
			Seq:    uint32(ls.sa.GetSectorTick()),
			Player: pEntity.ToPublicFB(),
		}})
}

func (ls *LoginSystem) handleDisconnect(ctx context.Context, sid float64) {
	log.Info().Str("SID", core.SIDStr(sid)).Msg("Player logout!")
	en, ok := ls.sidsToEnts[sid]
	if !ok {
		log.Error().Str("SID", core.SIDStr(sid)).Msg("Could not find entity for session during logout")
		return
	}

	ls.sa.RemoveEntity(en.ID())
	delete(ls.sidsToEnts, en.Sid)
}

func (ls *LoginSystem) handleLogout(ctx context.Context, sid float64) {
	en, ok := ls.sidsToEnts[sid]
	if !ok {
		log.Error().Str("SID", core.SIDStr(sid)).Msg("Could not find entity for session during logout")
		return
	}

	ls.handleDisconnect(ctx, sid)

	eb.Publish(eb.Event{
		Topic: eb.N_LOGOUT_RESPONSE,
		Data: eb.N_LOGOUT_RESPONSE_T{
			Sid: en.Sid,
			Seq: uint32(ls.sa.GetSectorTick()),
		}})
}
