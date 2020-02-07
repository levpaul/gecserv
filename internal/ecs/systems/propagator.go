package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs/components"
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
	})
}
func (pm *PropagatorSystem) Update(ctx context.Context, dt core.GameTick) {
	// loop over players
	// read their last state ack
	// calculate current interest zone and prev interest zone
	// get entities from current interest zones with changes
	// compare diff of old/new interest zones
	// publish diff as net_event
}
func (pm *PropagatorSystem) Add(en core.Entity)      {}
func (pm *PropagatorSystem) Remove(en core.EntityID) {}
