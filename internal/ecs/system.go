package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
)

// System is a base interface for Systems. This type of System will not operate directly
// on entities, but instead will likely interact with the event bus or other globals.
type System interface {
	Update(ctx context.Context, delta core.GameTick)
	SetSectorID(sectorID core.SectorID)
	GetSectorID() core.SectorID
}

// EntitySystem runs updates on a set of entities. Entities are added to systems in a sector
// automatically by the sectorAdmin. When a EntitySystem is first added to a sectorAdmin, an
// interface must be specified for the types of entities a given EntitySystem operates on.
type EntitySystem interface {
	System
	Add(en EntityID)
	Remove(en EntityID)
}

// Initializer is an interface that sectorAdmin checks for on any System passed to it and
// if implemented, the sectorAdmin with initialize via the Init() func.
type Initializer interface {
	Init()
}
