package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/levpaul/idolscape-backend/internal/ecs"
	"github.com/levpaul/idolscape-backend/internal/ecs/entities"
	"github.com/levpaul/idolscape-backend/internal/fb"
	"github.com/rs/zerolog/log"
)

// LoginSystem handles when a player has logged into a character for a
// given sector
type LoginSystem struct {
	BaseSystem
	logins chan eb.Event
}

// TODO: The topic for logins may need to be split per sector
func (ls *LoginSystem) Init() {
	ls.logins = make(chan eb.Event, 128)
	eb.Subscribe(eb.S_LOGIN, ls.logins)
}

func (ls *LoginSystem) Update(ctx context.Context, dt core.GameTick) {
	for {
		select {
		case l := <-ls.logins:
			player, ok := l.Data.(eb.S_LOGIN_T)
			if !ok {
				log.Error().Interface("data", l.Data).Msg("Failed to type assert S_LOGIN message")
				continue
			}
			ls.handleLogin(ctx, player)
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func (ls *LoginSystem) handleLogin(ctx context.Context, player *fb.PlayerT) {
	log.Info().Float64("SID", player.Sid).Msg("New player login!")

	pEntity := &entities.PlayerE{
		PlayerT: player,
	}

	err := ecs.AddEntityToSector(ls.sectorID, pEntity)
	if err != nil {
		log.Err(err).Msg("Failed to add player entity")
		return
	}

	// Steps:
	// 1 - Add player to map
	//     -> Create 'player' entity & push through sectorAdmin?
	//     -> Some sort of mapping store
	//     ->
	//     ->
	//     ->
	// 2 - Propagate to other players
	//     -> Post gametick propagation + interest system
	//     ->
	//     ->

}
