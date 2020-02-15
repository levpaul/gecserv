package systems

import (
	"context"
	"github.com/levpaul/gecserv/internal/core"
)

// PlayerMovementSystem is repsonsible for reading current inputs and moving player
// and updating interest mapping
type PlayerMovementSystem struct {
	BaseSystem
}

func (pm *PlayerMovementSystem) Init() {}
func (pm *PlayerMovementSystem) Update(ctx context.Context, dt core.GameTick) {
	// Loop through all players
	// If single input, process it
	// Else if multi input replay all buffered movements
	// Else if empty input, replay last input
}
