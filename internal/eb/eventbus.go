package eb

import (
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

const (
	publishTimeout = time.Second * 2
)

var (
	pipeErr chan<- error
	eb      eventBus
)

func Start(pErr chan<- error) error {
	pipeErr = pErr
	inititializeEventBus()
	return nil
}

type Event struct {
	Topic EventTopic
	Data  interface{}
}

type eventBus struct {
	subs map[EventTopic]subscribers
	rw   sync.RWMutex
}

type subscribers []chan<- Event

func Publish(ev Event) {
	eb.publish(ev)
}

func (e *eventBus) publish(ev Event) {
	e.rw.RLock()
	defer eb.rw.RUnlock()
	subCp := append(subscribers{}, e.subs[ev.Topic]...)
	go func(ev Event, subs subscribers) {
		timeout := time.After(publishTimeout)
		for _, ch := range subs {
			select {
			case ch <- ev:
			case <-timeout:
				log.Error().Interface("Event", ev).Msg("Failed to publish event")
				return
			}

		}
	}(ev, subCp)
}

func Subscribe(t EventTopic, subCh chan<- Event) {
	eb.subscribe(t, subCh)
}

func (e *eventBus) subscribe(t EventTopic, subCh chan<- Event) {
	e.rw.Lock()
	defer e.rw.Unlock()

	eb.subs[t] = append(eb.subs[t], subCh)
}
