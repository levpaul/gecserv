package ecs

import (
	"context"
	"github.com/levpaul/gecserv/internal/config"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
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
