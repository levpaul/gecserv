package eventbus

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

type eventBus struct {
	subs map[EventTopic]subscribers
	rw   sync.RWMutex
}

type subscribers []chan<- interface{}

func Publish(t EventTopic, data interface{}) {
	eb.publish(t, data)
}

func (e *eventBus) publish(t EventTopic, data interface{}) {
	e.rw.RLock()
	defer eb.rw.RUnlock()
	subCp := append(subscribers{}, e.subs[t]...)
	go func(data interface{}, subs subscribers) {
		timeout := time.After(publishTimeout)
		for _, ch := range subs {
			select {
			case ch <- data:
			case <-timeout:
				log.Error().Interface("Data", data).Int("Topic", int(t)).Msg("Failed to publish event")
				return
			}

		}
	}(data, subCp)
}

func Subscribe(t EventTopic, subCh chan<- interface{}) {
	eb.subscribe(t, subCh)
}

func (e *eventBus) subscribe(t EventTopic, subCh chan<- interface{}) {
	e.rw.Lock()
	defer e.rw.Unlock()

	eb.subs[t] = append(eb.subs[t], subCh)
}
