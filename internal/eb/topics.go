package eb

const (
	PLOGINS = iota
	PLOGOUTS
	PINPUT
	NCONNECT
	NDISCONN
)

type EventTopic int

func inititializeEventBus() {
	eb = eventBus{}
	eb.subs = make(map[EventTopic]subscribers)

	// Pre initialize all sub channels
	eb.subs[PLOGINS] = subscribers{}
	eb.subs[PLOGOUTS] = subscribers{}
	eb.subs[PINPUT] = subscribers{}
	eb.subs[NCONNECT] = subscribers{}
	eb.subs[NDISCONN] = subscribers{}
}
