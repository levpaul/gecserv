package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/rs/zerolog/log"
	"reflect"
	"sync"
)

var (
	sectorIDCounter   core.SectorID
	sectorIDCounterMu sync.Mutex
)

type sectorAdmin struct {
	id       core.SectorID
	systems  []System
	entities map[core.EntityID]Entity

	// used to automatically adding new entities to relevant systems
	entitySystemTypes map[reflect.Type]reflect.Type

	mu sync.Mutex
}

func newSectorAdmin() *sectorAdmin {
	sectorIDCounterMu.Lock()
	defer sectorIDCounterMu.Unlock()
	sectorIDCounter += 1
	sa := new(sectorAdmin)
	sa.id = sectorIDCounter
	sa.entities = make(map[core.EntityID]Entity)
	sa.entitySystemTypes = make(map[reflect.Type]reflect.Type)
	return sa
}

// Expects a pointer to a system, initializing if possible
func (sa *sectorAdmin) addSystem(s System) {
	if sysInit, ok := s.(Initializer); ok {
		sysInit.Init()
	}

	s.SetSectorID(sa.id)

	sa.systems = append(sa.systems, s)
}

func (sa *sectorAdmin) addEntitySystem(s EntitySystem, ifce interface{}) {
	sa.addSystem(s)

	// Create entry for system for given interface - many systems could
	// be used for a single entity type
	sa.entitySystemTypes[reflect.TypeOf(s)] = reflect.TypeOf(ifce)
}

func (sa *sectorAdmin) addEntity(en Entity) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	sa.entities[en.ID()] = en

	for _, s := range sa.systems {
		es, ok := s.(EntitySystem)
		if !ok {
			continue
		}

		if reflect.TypeOf(en) == sa.entitySystemTypes[reflect.TypeOf(es)] {
			es.Add(en.ID())
		}
	}
}

func (sa *sectorAdmin) removeEntity(en core.EntityID) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	for _, s := range sa.systems {
		if es, ok := s.(EntitySystem); ok {
			es.Remove(en)
		}
	}

	delete(sa.entities, en)
}

func (sa *sectorAdmin) update(ctx context.Context, dt core.GameTick) {
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

func (sa *sectorAdmin) getEntity(entityID core.EntityID) Entity {
	en, ok := sa.entities[entityID]
	if !ok {
		return nil
	}
	return en
}
