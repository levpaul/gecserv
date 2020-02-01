package main

import (
	"github.com/levpaul/idolscape-backend/internal/cmdflags"
	"github.com/levpaul/idolscape-backend/internal/config"
	"github.com/levpaul/idolscape-backend/internal/debug"
	"github.com/levpaul/idolscape-backend/internal/eb"
	"github.com/levpaul/idolscape-backend/internal/ecs"
	"github.com/levpaul/idolscape-backend/internal/ingest"
	"github.com/levpaul/idolscape-backend/internal/netconn"
	"github.com/levpaul/idolscape-backend/internal/propagator"
	"github.com/levpaul/idolscape-backend/internal/simulator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
		viper.SetConfigName("dev")
		go debug.StartDebugServer()
	}

	config.Init()

	startPipeline("eventbus", eb.Start)           // Manage message sharing channels
	startPipeline("ecs", ecs.Start)               // Manage game state
	startPipeline("netconn", netconn.Start)       // Manage data connections to clients
	startPipeline("simulator", simulator.Start)   // Simulate each game tick for each client
	startPipeline("propagator", propagator.Start) // Send updates to each client
	startPipeline("ingest", ingest.Start)         // Take client input + handle registration

	select {
	case err := <-pipelineErrCh:
		log.Err(err).Send()
		os.Exit(1)
	}
}

func startPipeline(plName string, pipeline func(chan<- error) error) {
	if err := pipeline(pipelineErrCh); err != nil {
		log.Err(err).Str("pipeline name", plName).Msg("Failed to start pipeline")
		os.Exit(1)
	}
}
