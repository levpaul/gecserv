package systems

import (
	"context"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/ecs/components"
	"github.com/levpaul/gecserv/internal/ecs/entities"
	"github.com/rs/zerolog/log"
)

const (
	MaxTickDiff = 15
)

// PropagatorSystem is repsonsible for reading sending relevant updates to players
type PropagatorSystem struct {
	BaseSystem
	cc core.ComponentCollection
}

func (pm *PropagatorSystem) Init() {
	pm.cc = core.NewComponentCollection([]interface{}{
		new(components.StateHistoryComponent),
		new(components.NetworkedSessionComponent),
		new(components.PositionalComponent),
	})
}
func (pm *PropagatorSystem) Update(ctx context.Context, dt core.GameTick) {
	interestMapEnts := pm.sa.FilterEntitiesByCC(core.NewComponentCollection([]interface{}{
		new(components.InterestMapComponent),
	}))
	ime := interestMapEnts.Next()
	if ime == nil || interestMapEnts.Next() != nil {
		log.Fatal().Msg("Unexpected interestmap exception, either interest map not found or more than 1 was!")
	}
	im := ime.(components.InterestMapComponent).GetInterestMap()

	ents := pm.sa.FilterEntitiesByCC(pm.cc)
	log.Info().Msgf("Ents", ents)
	for en := ents.Next(); en != nil; en = ents.Next() {
		log.Info().Msgf("En", en)
		entStateHist := en.(components.StateHistoryComponent).GetStateHistory()

		// Send full player list if lastAck is too far off
		if pm.sa.GetSectorTick()-entStateHist.LastAck > MaxTickDiff || entStateHist.LastAck == 0 {
			pm.sendFullState(en, im)
			continue
		}

		log.Warn().Msg("Unsupported partial diffs currently")
		// impl plan below VVVVV
		// lookup state from lastAck
		// get position from oldstate
		// determine old interestzone
		// determine new interestzone
		// for overlapping sectors send diffs
		// for removed sectors delete ??
		// push as net_events
	}
}

func (pm *PropagatorSystem) sendFullState(en core.Entity, im components.InterestMap) {
	log.Info().Uint32("Player Ent ID", uint32(en.ID())).Msg("Sending full state")

	// send logins
	// send logouts
	// send positions

	players := []*entities.PlayerE{}
	for i := range im.Imap {
		for j := range im.Imap[i] {
			for _, e := range im.Imap[i][j] {
				imEn := pm.sa.GetEntity(e)
				if imEn == nil {
					continue
				}
				plEn, ok := imEn.(*entities.PlayerE)
				if !ok {
					continue
				}
				players = append(players, plEn)
			}
		}
	}

	tp, ok := en.(*entities.PlayerE)
	if !ok {
		log.Fatal().Msg("Some strange shit happened")
	}

	log.Info().Msgf("Player list: %s", players)
	eb.Publish(eb.Event{
		Topic: eb.N_PLAYER_SYNC,
		Data: eb.N_PLAYER_SYNC_T(&eb.PlayerSyncMessage{
			ToPlayerSID: tp.Sid,
			Players:     players,
			Tick:        pm.sa.GetSectorTick(),
		})})
}
