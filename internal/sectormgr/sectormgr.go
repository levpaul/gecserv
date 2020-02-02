package sectormgr

import (
	"github.com/levpaul/idolscape-backend/internal/ecs"
	"github.com/levpaul/idolscape-backend/internal/ecs/entities"
	"github.com/levpaul/idolscape-backend/internal/ecs/systems"
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
	sectorID := ecs.AddNewSector()
	ecs.AddSystemToSector(new(systems.LoginSystem), sectorID)
	ecs.AddEntitySystemToSector(new(systems.LogoutSystem), &entities.PlayerE{}, sectorID)
}
