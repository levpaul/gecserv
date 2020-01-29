package simulator

import (
	"github.com/levpaul/idolscape-backend/internal/config"
	"sync"
	"time"
)

var (
	tickDoneListeners []chan struct{}
	tickDone          = make(chan struct{})
	tickDoneLock      sync.Mutex

	pipeErr chan<- error
)

func Start(pErr chan<- error) error {
	pipeErr = pErr
	go startTickDoneNotifier()
	go startSimulator()
	return nil
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
			tickDone <- struct{}{}
		}
	}
}

func simulate() error {

	//log.Info().Msg("Game tick!")

	// Get current list of chars from GAMESTATE

	// Check event bus for actions from current chars -> simulate/predict

	// Get logins/logoffs from event bus -> update GAMESTATE chars

	// end sim tick

	return nil
}

func GetTickDoneCh() <-chan struct{} {
	tickDoneLock.Lock()
	defer tickDoneLock.Unlock()
	tl := make(chan struct{})
	tickDoneListeners = append(tickDoneListeners, tl)
	return nil
}

// TODO: Revist this and make non-blocking
func startTickDoneNotifier() {
	for {
		<-tickDone
		for _, lch := range tickDoneListeners {
			lch <- struct{}{}
		}
	}
}
