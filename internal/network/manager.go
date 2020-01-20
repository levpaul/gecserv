package network

import (
	"fmt"
	"time"
)

const gameTick = time.Second / 20

func StartNetworkManager() {

	gameTicker := time.NewTicker(gameTick)

	for range gameTicker.C {
		// Send updates for new connections
		newConnLock.Lock()
		for _, c := range newConns {
			x, y := c.GS.GetPos()
			col := c.GS.GetCol()
			fmt.Println("new connection: ", x, y)
			for _, existC := range conns {
				if c.Uuid != existC.Uuid {
					fmt.Println("Sending update to a palyer", existC.Uuid)
					existC.SendNewPlayer(col, x, y)
				}
			}
		}
		newConns = nil
		newConnLock.Unlock()

		// Send position updates for everyone

	}
}
