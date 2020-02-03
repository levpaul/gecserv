package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs"
)

// InterestSystem is repsonsible for updating a singleton map of interest buckets
// containing all entities in subdivisions of the sectors map, used by the propagator
// to send relevant map state only to clients. InterestSystem listens for objectMove
// updates from the eventbus and updates all entities from there. May add a scheduled
// full update in too
type InterestSystem struct {
	BaseSystem
}

func (is *InterestSystem) Init()                                        {}
func (is *InterestSystem) Update(ctx context.Context, dt core.GameTick) {}
func (is *InterestSystem) Add(en ecs.Entity)                            {}
func (is *InterestSystem) Remove(en core.EntityID)                      {}
