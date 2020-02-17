package systems

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/ecs/components"
	"github.com/levpaul/gecserv/internal/ecs/entities"
	"github.com/levpaul/gecserv/internal/fb"
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
		new(components.NetworkSessionComponent),
		new(components.PositionComponent),
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
	for en := ents.Next(); en != nil; en = ents.Next() {
		shCp := en.(components.StateHistoryComponent).GetStateHistory()
		// read their last state ack
		if pm.sa.GetSectorTick()-shCp.LastAck > MaxTickDiff || shCp.LastAck == 0 {
			pm.sendAllPlayers(en, im)
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

func (pm *PropagatorSystem) sendAllPlayers(en core.Entity, im components.InterestMap) {
	log.Info().Uint32("Player", uint32(en.ID())).Msg("Sending full state")
	//players := []*entities.PlayerE{}
	players := []*fb.PlayerT{}
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
				players = append(players, plEn.ToFB())
			}
		}
	}
	// full player list here
	fmt.Println("Plkayer list: ", players)
	eb.Publish(eb.Event{
		Topic: eb.N_PLAYER_SYNC,
		Data: eb.N_PLAYER_SYNC_T(&eb.PlayerSyncMessage{
			ToPlayer: uuid.UUID{},
			Players:  players,
		})})
}
