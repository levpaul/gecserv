package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/config"
	"github.com/levpaul/idolscape-backend/internal/core"
	"time"
)

func updateLoop() {
	gt := time.Tick(config.GameTickDuration)

	for {
		select {
		case <-gt:
			if err := simulate(); err != nil {
				pipeErr <- err
				return
			}
			// TODO: Publish an EB message like "S_COMPLETE"
			//   - Then make propagator listen for that to cut diffs and send
		}
	}
}

func simulate() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.GameTickDuration)
	defer cancel()

	// This function should actually run through all the systems on a mapSegment, calling the updates
	// Segments should use priority such that an ordering like this occurs:

	// LoginSystem (adds player entities to map)
	// LogoutSystem (removes player entities from map)
	// PlayerMovementSystem
	// NPCMovementSystem
	// ...
	// Then?
	// ENDGAMETICK -> Publish Prop msg -> Begin propagation system update
	// So maybe have two main "Worlds" Sim/Prop?

	// Each system can read/push from the message bus to during their update calls

updateLoop:
	for _, s := range sectors {
		select {
		case <-ctx.Done():
			break updateLoop
		default:
		}
		s.update(ctx, core.GameTick(1))
	}
	return nil
}
