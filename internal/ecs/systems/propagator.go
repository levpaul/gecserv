package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs/components"
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
	ents := pm.sa.FilterEntitiesByCC(pm.cc)
	for en := ents.Next(); en != nil; en = ents.Next() {
		shCp := en.(components.StateHistoryComponent).GetStateHistory()
		// read their last state ack
		if pm.sa.GetSectorTick()-shCp.LastAck > MaxTickDiff || shCp.LastAck == 0 {
			pm.sendFullState(en)
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

func (pm *PropagatorSystem) sendFullState(en core.Entity) {
	log.Info().Msg("Sending full state")
}
