package network

import (
	"github.com/rs/zerolog/log"
	"time"
)

const gameTick = time.Second / 20

func StartNetworkManager() {

	gameTicker := time.NewTicker(gameTick)

	// Update Loop
	for range gameTicker.C {
		// Send updates for new connections
		newConnLock.Lock()
		for _, c := range newConns {
			log.Info().Interface("Conn", c).Msg("New connection")
			for _, liveC := range conns {
				if c.Uuid != liveC.Uuid {
					log.Info().Str("UUID", liveC.Uuid.String()).Msg("Sending update to player")
					liveC.SendNewPlayer(c)
				}
			}
		}
		newConns = nil
		newConnLock.Unlock()

		// Send disconnects to everyone
		discConnLock.Lock()
		for _, c := range discConns {
			for _, liveC := range conns {
				liveC.SendDisconnectedPlayer(c)
			}
		}
		discConns = nil
		discConnLock.Unlock()
	}

	panic("Network manager game loop exited!")
}
