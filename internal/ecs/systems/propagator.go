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
	for en := ents.Next(); en != nil; en = ents.Next() {
		entStateHist := en.(components.StateHistoryComponent).GetStateHistory()

		// Send full player list if lastAck is too far off
		if pm.sa.GetSectorTick()-entStateHist.LastAck > MaxTickDiff || entStateHist.LastAck == 0 {
			pm.sendCurrentFullState(en, im)
			continue
		}

		log.Warn().Msg("Unsupported partial diffs currently")
		// lookup state from lastAck
		// get position from oldstate
		// determine old interestzone
		// determine new interestzone
		// for overlapping sectors send diffs
		// for removed sectors delete ??
		// push as net_events
	}
}

func (pm *PropagatorSystem) sendCurrentFullState(en core.Entity, im components.InterestMap) {
	currPlayerEn, ok := en.(*entities.PlayerE)
	if !ok {
		log.Fatal().Msg("Some strange shit happened - entity is not a player??")
	}

	currTick := pm.sa.GetSectorTick()

	players := []*fb.PlayerT{}
	logins := []*fb.PlayerT{}
	logouts := []float64{}

	for i := range im.Imap {
		for j := range im.Imap[i] {
			for _, e := range im.Imap[i][j] {
				imEn := pm.sa.GetEntity(e)
				if imEn == nil {
					continue
				}
				plEn, ok := imEn.(*entities.PlayerE)
				if !ok { // Entity is not a player
					continue
				}
				if plEn.LoginTick == currTick {
					logins = append(logins, plEn.ToPublicFB())
				} else {
					players = append(players, plEn.ToPublicFB())
				}
			}
		}
	}

	eb.Publish(eb.Event{
		Topic: eb.N_PLAYER_SYNC,
		Data: eb.N_PLAYER_SYNC_T{
			ToPlayerSID: currPlayerEn.Sid,
			Msg: fb.MapUpdateT{
				Seq:     uint32(currTick),
				Logins:  logins,
				Logouts: logouts,
				Psyncs:  players,
			},
		}})
}
