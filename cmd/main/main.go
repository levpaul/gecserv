package main

import (
	"errors"
	"github.com/levpaul/idolscape-backend/internal/cmdflags"
	"github.com/levpaul/idolscape-backend/internal/debug"
	"github.com/levpaul/idolscape-backend/internal/flusher"
	"github.com/levpaul/idolscape-backend/internal/ingest"
	"github.com/levpaul/idolscape-backend/internal/network"
	"github.com/levpaul/idolscape-backend/internal/propagation"
	"github.com/levpaul/idolscape-backend/internal/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	cmdflags.Parse()
}

var pipelineErrCh = make(chan error)

func main() {
	if *cmdflags.DevMode == true {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		go debug.StartDebugServer()
	}

	// Legacy - to be replaced by flusher pipeline
	go network.StartNetworkManager()

	startPipeline("ingest", ingest.Start)
	startPipeline("validation", validation.Start)
	startPipeline("flusher", flusher.Start)
	startPipeline("propagation", propagation.Start)

	select {
	case err := <-pipelineErrCh:
		log.Err(err).Send()
		return
	}
}

func startPipeline(plName string, pipeline func(chan<- error)) {
	go func() {
		pipeline(pipelineErrCh)
		pipelineErrCh <- errors.New("pipeline job returned unexpectedly")
	}()
}
