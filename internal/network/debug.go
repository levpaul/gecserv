package network

import (
	"errors"
	"github.com/levpaul/idolscape-backend/internal/cmdflags"
	"github.com/rs/zerolog/log"
)

func DebugGetLiveConnection() (Connection, error) {
	if *cmdflags.DevMode == false {
		log.Fatal().Msg("Debug function called without Devmode enabled")
	}
	if len(conns) == 0 {
		return Connection{}, errors.New("no active connections found")
	}
	return conns[0], nil
}
