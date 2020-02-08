package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs/entities"
	"github.com/rs/zerolog/log"
	"reflect"
	"sync"
)

var (
	sectorIDCounter   core.SectorID
	sectorIDCounterMu sync.Mutex
)

type SectorAdmin struct {
	id         core.SectorID
	sectorTick core.GameTick
	systems    []System
	entities   map[core.EntityID]core.Entity

	// Singletons
	sectorMap   *entities.MapE
	playerList  map[float64]*entities.PlayerE
	interestMap *[][]core.EntityIDs
}

func newSectorAdmin() *SectorAdmin {
	sectorIDCounterMu.Lock()
	defer sectorIDCounterMu.Unlock()

	sa := new(SectorAdmin)
	sectorIDCounter += 1
	sa.id = sectorIDCounter
	sa.entities = make(map[core.EntityID]core.Entity)

	sa.sectorMap = entities.NewDefaultMap()
	sa.AddEntity(sa.sectorMap)

	return sa
}

// Expects a pointer to a system, initializing if possible
func (sa *SectorAdmin) AddSystem(s System) {
	s.SetSectorAdmin(sa)
	if sysInit, ok := s.(Initializer); ok {
		sysInit.Init()
	}

	sa.systems = append(sa.systems, s)
}

func (sa *SectorAdmin) AddEntity(en core.Entity) {
	sa.entities[en.ID()] = en
}

func (sa *SectorAdmin) RemoveEntity(en core.EntityID) {
	delete(sa.entities, en)
}

type entityIterator struct {
	ents []core.Entity
	idx  int
}

func (ei *entityIterator) Next() core.Entity {
	if ei.idx >= len(ei.ents) {
		return nil
	}
	ei.idx += 1
	return ei.ents[ei.idx-1]
}

// TODO: This is a major area for optimisation!!!!!
// One idea might be to follow the rust/hecs pattern where entities are stored in arrays
// based on their entity type. Then you can check the first element of one of those arrays
// and if it matches you can just return the slice of all of those entites, since they all
// match the CC. In this case you would return a [][]core.Entity for the system to map over
// or maybe create a helper struct to act as an iterator. The benefit of this approach aims
// to try and make memory access efficient by returning contiguous runs of entities relevant
// to the systems. To do this in go though you need to control the allocation of the entity
// slices, which I'm not sure how feasible that is for Go - my worry is that an interface
// slice will just hold pointers to the underlying entity structs, which defeats the point
// of an attempted optimization.
func (sa *SectorAdmin) FilterEntitiesByCC(cc core.ComponentCollection) core.EntityIterator {
	reEn := []core.Entity{}
	for _, e := range sa.entities {
		satisfiesCC := true
		for _, t := range cc {
			if !reflect.TypeOf(e).Implements(t) {
				satisfiesCC = false
				break
			}
		}
		if satisfiesCC {
			reEn = append(reEn, e)
		}
	}
	return &entityIterator{
		ents: reEn,
		idx:  0,
	}
}

func (sa *SectorAdmin) GetEntity(entityID core.EntityID) core.Entity {
	en, ok := sa.entities[entityID]
	if !ok {
		return nil
	}
	return en
}

func (sa *SectorAdmin) update(ctx context.Context, dt core.GameTick) {
	sa.sectorTick += 1
	for _, s := range sa.systems {
		select {
		case <-ctx.Done():
			log.Warn().Uint32("SectorID", uint32(sa.id)).Msg("Context timout in sectoradmin update")
			return
		default:
			s.Update(ctx, dt)
		}
	}
}

// ============ SINGLETONS =============

func (sa *SectorAdmin) SetPlayerListSingleton(pl map[float64]*entities.PlayerE) {
	if sa.playerList != nil {
		log.Error().Msg("tried to set playerlist singleton which has already been set")
		return
	}
	sa.playerList = pl
}

// GetPlayerListSingleton returns a current map of session IDs to player entities -
// this should not be written to by any callers except for whatever called
// SetPlayerListSingleton
func (sa *SectorAdmin) GetPlayerListSingleton() map[float64]*entities.PlayerE {
	return sa.playerList
}

func (sa *SectorAdmin) SetInterestMapSingleton(im *[][]core.EntityIDs) {
	if sa.interestMap != nil {
		log.Error().Msg("tried to set interest map singleton which has already been set")
		return
	}
	sa.interestMap = im
}

// GetPlayerListSingleton returns a current map of session IDs to player entities -
// this should not be written to by any callers except for whatever called
// SetPlayerListSingleton
func (sa *SectorAdmin) GetInterestMapSingleton() [][]core.EntityIDs {
	return *sa.interestMap
}

func (sa *SectorAdmin) GetSectorTick() core.GameTick {
	return sa.sectorTick
}
