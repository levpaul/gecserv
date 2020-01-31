package simulator

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/config"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/rs/zerolog/log"
	"reflect"
	"time"
)

var (
	pipeErr chan<- error

	busCh chan eb.Event
)

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	go startSimulator()
	return nil
}

func initialize() {
	busCh = make(chan eb.Event, 128)

	eb.Subscribe(eb.S_LOGIN, busCh)
	eb.Subscribe(eb.S_LOGOUT, busCh)
}

func startSimulator() {
	// Start timer
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

	for {
		select {
		case e := <-busCh:
			switch e.Topic {
			case eb.S_LOGIN:
				//fpt := e.Data.(*fb.PlayerT)
				data, ok := e.Data.(eb.S_LOGIN_T)
				if !ok {
					log.Error().Interface("data", e.Data).Msg("Failed to type assert S_LOGIN message")
					log.Error().Interface("type", reflect.TypeOf(e.Data)).Send()
					continue
				}
				handleLogin(ctx, eb.S_LOGIN_T(data))
			case eb.S_LOGOUT:
				data, ok := e.Data.(eb.S_LOGOUT_T)
				if !ok {
					log.Error().Interface("data", e.Data).Msg("Failed to type assert S_LOGOUT message")
					continue
				}
				handleLogout(ctx, data)
			}

		case <-ctx.Done():
			log.Ctx(ctx).Info().Send()
		}
	}

	return nil
}

func handleLogin(ctx context.Context, e eb.S_LOGIN_T) {
	log.Info().Float64("SID", e.Sid).Msg("SOMEONE LOGIN")
	// TODO: Add player to map
}

func handleLogout(ctx context.Context, e eb.S_LOGOUT_T) {
	log.Info().Float64("SID", float64(e)).Msg("SOMEONE LOGout")
}
