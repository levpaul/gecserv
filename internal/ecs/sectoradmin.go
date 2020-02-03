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
	entities map[core.EntityID]Entity

	// used to automatically adding new entities to relevant systems
	entitySystemTypes map[reflect.Type]reflect.Type

	// Singletons
	sectorMap   *entities.MapE
	playerList  map[float64]*entities.PlayerE
	interestMap [][]Entity

	mu sync.Mutex
}

func newSectorAdmin() *SectorAdmin {
	sectorIDCounterMu.Lock()
	defer sectorIDCounterMu.Unlock()

	sa := new(SectorAdmin)
	sectorIDCounter += 1
	sa.id = sectorIDCounter
	sa.entities = make(map[core.EntityID]Entity)
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

func (sa *SectorAdmin) AddEntitySystem(s EntitySystem, ifce interface{}) {
	sa.AddSystem(s)

	// Create entry for system for given interface - many systems could
	// be used for a single entity type
	sa.entitySystemTypes[reflect.TypeOf(s)] = reflect.TypeOf(ifce)
}

func (sa *SectorAdmin) AddEntity(en Entity) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	sa.entities[en.ID()] = en

	for _, s := range sa.systems {
		es, ok := s.(EntitySystem)
		if !ok {
			continue
		}

		if reflect.TypeOf(en) == sa.entitySystemTypes[reflect.TypeOf(es)] {
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

	delete(sa.entities, en)
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

func (sa *SectorAdmin) GetEntity(entityID core.EntityID) Entity {
	en, ok := sa.entities[entityID]
	if !ok {
		return nil
	}
	return en
}

// ============ SINGLETONS =============

// This file contains all global getters for singleton entities
func (sa *SectorAdmin) SetPlayerList(pl map[float64]*entities.PlayerE) {
	if sa.playerList != nil {
		log.Error().Msg("tried to set playerlist singleton which has already been set")
		return
	}
	sa.playerList = pl
}

// GetPlayerList returns a current map of session IDs to player entities -
// this should not be written to by any callers except for whatever called
// SetPlayerList
func (sa *SectorAdmin) GetPlayerList() map[float64]*entities.PlayerE {
	return sa.playerList
}
