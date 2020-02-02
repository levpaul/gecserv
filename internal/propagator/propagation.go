package propagator

var pipeErr chan<- error

func Start(pErr chan<- error) error {
	pipeErr = pErr
	go start()
	return nil
}

func start() {
	for {
		select {
		// Listen to propagate event and go through all players, sending updates
		}
	}
}
