package state

import (
	"github.com/levpaul/idolscape-backend/internal/fb"
)

var pipeErr chan<- error

var mapState fb.MapT

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	go start()
	return nil
}

func initialize() {
	mapState = fb.MapT{}
}

func start() {
	for {
		select {}
	}
}

func AddPlayer() {}
