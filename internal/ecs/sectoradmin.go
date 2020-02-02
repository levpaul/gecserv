package ecs

import (
	"context"
	"fmt"
	"github.com/levpaul/idolscape-backend/internal/core"
	"reflect"
	"sync"
)

type sectorAdmin struct {
	id       core.SectorID
	systems  []System
	entities map[EntityID]Entity

	// used to automatically adding new entities to relevant systems
	entitySystemInterfaces map[reflect.Type]reflect.Type

	mu sync.Mutex
}

func newSectorAdmin(SID core.SectorID) *sectorAdmin {
	sa := new(sectorAdmin)
	sa.id = SID
	sa.entities = make(map[EntityID]Entity)
	sa.entitySystemInterfaces = make(map[reflect.Type]reflect.Type)
	return sa
}

// Expects a pointer to a system, initializing if possible
func (sa *sectorAdmin) addSystem(s System) {
	if sysInit, ok := s.(Initializer); ok {
		sysInit.Init()
	}
	sa.systems = append(sa.systems, s)
}

func (sa *sectorAdmin) addEntitySystem(s EntitySystem, ifce interface{}) {
	sa.addSystem(s)

	// Create entry for system for given interface - many systems could
	// be used for a single entity type
	sa.entitySystemInterfaces[reflect.TypeOf(s)] = reflect.TypeOf(ifce)
}

func (sa *sectorAdmin) addEntity(en Entity) {
	for _, s := range sa.systems {
		fmt.Println(s)
	}
}

func (sa *sectorAdmin) removeEntity(en EntityID) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	for _, s := range sa.systems {
		if es, ok := s.(EntitySystem); ok {
			es.Remove(en)
		}
	}
}

func (sa *sectorAdmin) Update(ctx context.Context, dt core.GameTick) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	for _, s := range sa.systems {
		s.Update(ctx, dt)
	}
}
