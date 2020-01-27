package connections

var pipeErr chan<- error

func Start(pErr chan<- error) error {
	pipeErr = pErr
	go start()
	return nil
}

func start() {
	for {
		select {}
	}
}

// StoreConn
// GetConn
// DelConn
