package propagator

import "github.com/levpaul/idolscape-backend/internal/eb"

var (
	pipeErr          chan<- error
	gameTickListener chan eb.Event
)

func Start(pErr chan<- error) error {
	pipeErr = pErr

	gameTickListener = make(chan eb.Event)
	eb.Subscribe(eb.S_GAMETICK_DONE, gameTickListener)

	go start()
	return nil
}

func start() {
	for {
		select {
		case <-gameTickListener:
			propagate()
		}
	}
}

func propagate() {
	// Non-System Way
	// Iterate through every player of every sector
	// Find interest zone
	// Cut diff between their last ack seq and current
	// Send diff

	// System way would be the same except it's just an EntitySystem
	// that works off PlayerE types?

	// Would be good to reduce EB usage
	// Network load should be spread as wide as possible though,
	// it shouldn't block next ticks not be limited to half a tick?

	// Maybe best plan is to start as a separate EB driven service
	// and see what the bandwidth/CPU limits are, and decide to bring
	// it in to ECS loop if better suited. The only extra work of EB approach
	// is that there will need to be a system capturing active players
	// for the propagator to know about active pEntities
}
