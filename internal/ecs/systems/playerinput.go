package systems

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
)

// PlayerInputSystem is repsonsible for reading current inputs and moving player
// and updating interest mapping
type PlayerInputSystem struct {
	BaseSystem
}

func (pm *PlayerInputSystem) Init() {}
func (pm *PlayerInputSystem) Update(ctx context.Context, dt core.GameTick) {
	// Loop through player input messages
	// add to player input buffers
}
func (pm *PlayerInputSystem) Add(en core.Entity)      {}
func (pm *PlayerInputSystem) Remove(en core.EntityID) {}
