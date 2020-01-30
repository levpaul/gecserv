package eb

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

type EventTopic int

func inititializeEventBus() {
	eb = eventBus{}
	eb.subs = make(map[EventTopic]subscribers)

	// Initialize all topics
	for t := EventTopic(0); t < NUMTOPICS; t += 1 {
		eb.subs[t] = subscribers{}
	}
}
