package simulator

//
//const gameTick = time.Second / 20
//
//func StartNetworkManager() {
//
//	gameTicker := time.NewTicker(gameTick)
//
//	// Update Loop
//	for range gameTicker.C {
//		// Send updates for new connections
//		ingest.newConnLock.Lock()
//		for _, c := range ingest.newConns {
//			log.Info().Interface("Conn", c).Msg("New connection")
//			for _, liveC := range ingest.conns {
//				if c.Uuid != liveC.Uuid {
//					log.Info().Str("UUID", liveC.Uuid.String()).Msg("Sending update to player")
//					liveC.SendNewPlayer(c)
//				}
//			}
//		}
//		ingest.newConns = nil
//		ingest.newConnLock.Unlock()
//
//		// Send disconnects to everyone
//		ingest.discConnLock.Lock()
//		for _, c := range ingest.discConns {
//			for _, liveC := range ingest.conns {
//				liveC.SendDisconnectedPlayer(c)
//			}
//		}
//		ingest.discConns = nil
//		ingest.discConnLock.Unlock()
//	}
//
//	panic("Network manager game loop exited!")
//}
