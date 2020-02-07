package systems

import (
	"context"
	"fmt"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs/components"
	"github.com/levpaul/idolscape-backend/internal/ecs/entities"
	"github.com/rs/zerolog/log"
)

// InterestSystem is responsible for updating a singleton map of interest buckets
// containing all entities in subdivisions of the sectors map, used by the propagator
// to send relevant map state only to clients. InterestSystem listens for objectMove
// updates from the eventbus and updates all entities from there. May add a scheduled
// full update in too
type InterestSystem struct {
	BaseSystem
	interestMap [][]entities.InterestZone
	cc          core.ComponentCollection
}

func (is *InterestSystem) Init() {
	is.interestMap = [][]entities.InterestZone{}
	is.sa.SetInterestMapSingleton(&is.interestMap)
	is.cc = core.NewComponentCollection([]interface{}{
		new(components.ChangeableComponent),
		new(components.PositionComponent),
	})
}
func (is *InterestSystem) Update(ctx context.Context, dt core.GameTick) {
	// Loop through all changeable entities w/ position
	// If changed, update interest map w/ new coordinates

	for _, en := range is.sa.FilterEntitiesByCC(is.cc) {
		chEn, ok := en.(components.ChangeableComponent)
		if !ok || !chEn.GetChangeable().Changed {
			continue
		}
		log.Info().Msg("I'm supposed to update the ent interest map here")
		// Magically updates interestMap
		chEn.GetChangeable().Changed = false
	}
}

func (is *InterestSystem) Add(en core.Entity) {
	chEn := en.(components.ChangeableComponent).GetChangeable()
	fmt.Println(chEn.Changed)
	chEn.Changed = true
}

func (is *InterestSystem) Remove(en core.EntityID) {
	e := is.sa.GetEntity(en)
	chEn := e.(components.ChangeableComponent).GetChangeable()
	fmt.Println(chEn.Changed)
	chEn.Changed = true
}
