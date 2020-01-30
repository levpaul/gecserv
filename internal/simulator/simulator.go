package simulator

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/config"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	pipeErr chan<- error

	loginCh  chan eb.Event
	logoutCh chan eb.Event
)

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	go startSimulator()
	return nil
}

func initialize() {
	loginCh = make(chan eb.Event, 128)
	logoutCh = make(chan eb.Event, 128)

	eb.Subscribe(eb.S_LOGIN, loginCh)
	eb.Subscribe(eb.S_LOGOUT, logoutCh)
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
		}
	}
}

func simulate() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.GameTickDuration)
	defer cancel()

	go handleLogins(ctx)
	go handleLogouts(ctx)

	select {
	case <-ctx.Done():
		log.Ctx(ctx).Info().Send()
	}

	return nil
}

func handleLogins(ctx context.Context) {
	for {
		select {
		case <-loginCh:
			log.Info().Msg("SOMEONE LOGIN")
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func handleLogouts(ctx context.Context) {
	for {
		select {
		case <-logoutCh:
			log.Info().Msg("SOMEONE LOGOUT")
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}
