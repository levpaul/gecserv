package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
)

type SectorID uint16

type SectorAdmin struct {
	id       SectorID
	systems  []System
	entities map[EntityID]*Entity

	// From overwatch ECS:
	// object_pool []Component
	// array []Component
}

// Expects a pointer to a system, initializing if possible
func (sa *SectorAdmin) AddSystem(s System) {
	if sysInit, ok := s.(Initializer); ok {
		sysInit.Init()
	}
	sa.systems = append(sa.systems, s)
}

func (sa *SectorAdmin) Update(ctx context.Context, dt core.GameTick) {
	for _, s := range sa.systems {
		s.Update(ctx, dt)
	}
}
