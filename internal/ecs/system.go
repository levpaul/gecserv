package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
)

// System is a base interface for Systems. This type of System will not operate directly
// on entities, but instead will likely interact with the event bus or other globals.
type System interface {
	Update(ctx context.Context, delta core.GameTick)
	SetSectorAdmin(sa *SectorAdmin)
	GetSectorAdmin() *SectorAdmin
}

// EntitySystem runs updates on a set of entities. Entities are added to systems in a sector
// automatically by the SectorAdmin. When a EntitySystem is first added to a SectorAdmin, an
// interface must be specified for the types of entities a given EntitySystem operates on.
type EntitySystem interface {
	System
	Add(en core.EntityID)
	Remove(en core.EntityID)
}

// Initializer is an interface that SectorAdmin checks for on any System passed to it and
// if implemented, the SectorAdmin with initialize via the Init() func.
type Initializer interface {
	Init()
}
