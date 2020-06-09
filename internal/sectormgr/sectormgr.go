package sectormgr

import (
	"github.com/levpaul/gecserv/internal/ecs"
	"github.com/levpaul/gecserv/internal/ecs/systems"
)

var pipeErr chan<- error

func Start(pErr chan<- error) error {
	pipeErr = pErr
	start()
	return nil
}

// TODO: This mgr should load sectors from a DB and/or a sector queue for now
//  it just makes a single default sector though
func start() {
	addDefaultSector()
}

func addDefaultSector() {
	sa := ecs.AddNewSector()

	sa.AddSystem(new(systems.LoginSystem))
	sa.AddSystem(new(systems.PlayerSystem))
	sa.AddSystem(new(systems.InterestSystem))
	sa.AddSystem(new(systems.PropagatorSystem))
}
