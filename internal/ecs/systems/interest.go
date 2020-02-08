package systems

import (
	"context"
	"fmt"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs/components"
	"github.com/rs/zerolog/log"
)

const (
	segmentsX = 50
	segmentsY = 50
)

type imCoord struct {
	x uint8
	y uint8
}

// InterestSystem is responsible for updating a singleton map of interest buckets
// containing all entities in subdivisions of the sectors map, used by the propagator
// to send relevant map state only to clients. InterestSystem listens for objectMove
// updates from the eventbus and updates all entities from there. May add a scheduled
// full update in too
type InterestSystem struct {
	BaseSystem
	interestMap [][]core.EntityIDs
	imLookup    map[core.EntityID]imCoord
	cc          core.ComponentCollection
}

func (is *InterestSystem) Init() {
	is.interestMap = make([][]core.EntityIDs, segmentsX)
	for i := range is.interestMap {
		is.interestMap[i] = make([]core.EntityIDs, segmentsY)
	}
	is.imLookup = make(map[core.EntityID]imCoord)

	is.sa.SetInterestMapSingleton(&is.interestMap)
	is.cc = core.NewComponentCollection([]interface{}{
		new(components.ChangeableComponent),
		new(components.PositionComponent),
	})
}
func (is *InterestSystem) Update(ctx context.Context, dt core.GameTick) {
	// Loop through all changeable entities w/ position
	// If changed, update interest map w/ new coordinates

	ents := is.sa.FilterEntitiesByCC(is.cc)
	for en := ents.Next(); en != nil; en = ents.Next() {
		chCp, ok := en.(components.ChangeableComponent)
		if !ok || !chCp.GetChangeable().Changed {
			continue
		}

		eid := en.ID()
		posCp, ok := en.(components.PositionComponent)
		if !ok {
			log.Error().Uint32("entity", uint32(eid)).Msg("Failed to turn entity into position component at interest system")
			continue
		}
		imPosX := uint8(posCp.GetPosition().X / segmentsX)
		imPosY := uint8(posCp.GetPosition().Y / segmentsY)

		// Check for new entity in sector
		old, isInLookup := is.imLookup[eid]
		if !isInLookup {
			is.interestMap[imPosX][imPosY] = append(is.interestMap[imPosX][imPosY], en.ID())
			is.imLookup[eid] = imCoord{imPosX, imPosY}
			continue
		}

		// No sector position update, skip
		if old.x == imPosX && old.y == imPosY {
			continue
		}

		// Update sector in IM by deleting old entry and adding new
		for i, v := range is.interestMap[old.x][old.y] {
			if v == eid {
				is.interestMap[old.x][old.y][i] = is.interestMap[old.x][old.y][len(is.interestMap[old.x][old.y])-1]
				is.interestMap[old.x][old.y] = is.interestMap[old.x][old.y][:len(is.interestMap[old.x][old.y])-1]
				is.interestMap[imPosX][imPosY] = append(is.interestMap[imPosX][imPosY], eid)
				break
			}
		}
		chCp.GetChangeable().Changed = false
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
