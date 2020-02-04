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
	id       core.SectorID
	systems  []System
	entities map[core.EntityID]core.Entity

	// used to automatically adding new entities to relevant systems
	entitySystemTypes map[reflect.Type]reflect.Type

	// Singletons
	sectorMap   *entities.MapE
	playerList  map[float64]*entities.PlayerE
	interestMap *[][]core.Entity

	// Entity DLL state
	eHead core.Entity
	eTail core.Entity

	mu sync.Mutex
}

func newSectorAdmin() *SectorAdmin {
	sectorIDCounterMu.Lock()
	defer sectorIDCounterMu.Unlock()

	sa := new(SectorAdmin)
	sectorIDCounter += 1
	sa.id = sectorIDCounter
	sa.entities = make(map[core.EntityID]core.Entity)
	sa.entitySystemTypes = make(map[reflect.Type]reflect.Type)

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

// AddEntitySystem takes an EntitySystem and a _pointer_ to an interface type
// which the SectorAdmin will use to dynamically add and remove all future
// entities which map the passed type.
func (sa *SectorAdmin) AddEntitySystem(s EntitySystem, ifce interface{}) {
	sa.AddSystem(s)

	// Create entry for system for given interface - many systems could
	// be used for a single entity type
	sa.entitySystemTypes[reflect.TypeOf(s)] = reflect.TypeOf(ifce).Elem()
}

func (sa *SectorAdmin) AddEntity(en core.Entity) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	// Update DLL
	if sa.eHead == nil {
		sa.eHead = en
		sa.eTail = en
	} else {
		en.SetPrev(sa.eTail)
		sa.eTail.SetNext(en)
		sa.eTail = en
	}

	sa.entities[en.ID()] = en

	// Add to relevant EntitySystems where entity matches specified component interface
	for _, s := range sa.systems {
		es, ok := s.(EntitySystem)
		if !ok {
			continue
		}

		if reflect.TypeOf(en).Implements(sa.entitySystemTypes[reflect.TypeOf(es)]) {
			es.Add(en)
		}
	}
}

func (sa *SectorAdmin) RemoveEntity(en core.EntityID) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	for _, s := range sa.systems {
		if es, ok := s.(EntitySystem); ok {
			es.Remove(en)
		}
	}

	// Remove from entity DLL and Map
	switch {
	case len(sa.entities) == 1:
		sa.eHead = nil
		sa.eTail = nil
	case sa.eHead.ID() == en:
		sa.eHead = sa.entities[en].Next()
		sa.eHead.SetPrev(nil)
	case sa.eTail.ID() == en:
		sa.eTail = sa.entities[en].Prev()
		sa.eTail.SetNext(nil)
	default:
		sa.entities[en].Prev().SetNext(sa.entities[en].Next())
		sa.entities[en].Next().SetPrev(sa.entities[en].Prev())
	}
	delete(sa.entities, en)
}

func (sa *SectorAdmin) GetEntity(entityID core.EntityID) core.Entity {
	en, ok := sa.entities[entityID]
	if !ok {
		return nil
	}
	return en
}

func (sa *SectorAdmin) GetEntitiesHead() core.Entity {
	return sa.eHead
}

func (sa *SectorAdmin) update(ctx context.Context, dt core.GameTick) {
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

func (sa *SectorAdmin) SetInterestMapSingleton(im *[][]core.Entity) {
	if sa.interestMap != nil {
		log.Error().Msg("tried to set interest map singleton which has already been set")
		return
	}
	sa.interestMap = im
}

// GetPlayerListSingleton returns a current map of session IDs to player entities -
// this should not be written to by any callers except for whatever called
// SetPlayerListSingleton
func (sa *SectorAdmin) GetInterestMapSingleton() [][]core.Entity {
	return *sa.interestMap
}
