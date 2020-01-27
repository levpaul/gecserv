package events

var pipeErr chan<- error

func Start(pErr chan<- error) error {
	pipeErr = pErr
	// ==== attack example (hit) ====
	// - pull message from queue
	// - locks user
	// - checks usersLastSeq - if less, then done + unlock
	// -
	// - updates gameTick state for char
	// - unlock

	go start()
	return nil
}

func start() {
	for {
		select {}
	}
}
