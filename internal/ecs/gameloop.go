package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/config"
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/eb"
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
			eb.Publish(eb.Event{
				Topic: eb.S_GAMETICK_DONE,
				Data:  eb.S_GAMETICK_DONE_T{},
			})
			// TODO: Publish an EB message like "S_COMPLETE"
			//   - Then make propagator listen for that to cut diffs and send
		}
	}
}

func simulate() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.GameTickDuration)
	defer cancel()

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
