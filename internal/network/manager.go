package network

import (
	"fmt"
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
			fmt.Println("new connection: %-v", c)
			for _, liveC := range conns {
				if c.Uuid != liveC.Uuid {
					fmt.Println("Sending update to a player", liveC.Uuid)
					liveC.SendNewPlayer(c)
				}
			}
		}
		newConns = nil
		newConnLock.Unlock()

		// Send disconnects to everyone
		discConnLock.Lock()
		for _, c := range discConns {
			fmt.Println("Starting disconnect update")
			for _, liveC := range conns {
				fmt.Printf("Discon: %v %v \n", c.Uuid, liveC.Uuid)
				liveC.SendDisconnectedPlayer(c)
			}
		}
		discConns = nil
		discConnLock.Unlock()

	}
}
