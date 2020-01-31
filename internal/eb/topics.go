package eb

import (
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/fb"
)

const (
	// Simulation messages
	S_LOGIN = 0 + iota
	S_LOGOUT
	S_INPUT
	// Network messages
	N_CONNECT = 128 + iota
	N_DISCONN
	// Used as marker for event topics - insert new topics ABOVE this one
	NUMTOPICS = 255
)

// Bindings of message topics to types
type (
	S_LOGIN_T  *fb.PlayerT
	S_LOGOUT_T float64
	S_INPUT_T  int
	// Network messages
	N_CONNECT_T *NetworkConnection
	N_DISCONN_T *uuid2.UUID
)

type EventTopic int

func inititializeEventBus() {
	eb = eventBus{}
	eb.subs = make(map[EventTopic]subscribers)

	// Initialize all topics
	for t := EventTopic(0); t < NUMTOPICS; t += 1 {
		eb.subs[t] = subscribers{}
	}
}
