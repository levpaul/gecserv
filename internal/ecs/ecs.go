package ecs

import (
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs/systems"
)

var (
	pipeErr chan<- error

	sectors map[core.SectorID]*sectorAdmin
)

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	go updateLoop()
	return nil
}

func initialize() {
	sectors = map[core.SectorID]*sectorAdmin{
		1: newSectorAdmin(1),
	}

	// Load default systems for all sectors
	for _, s := range sectors {
		s.addSystem(new(systems.LoginSystem))
		s.addSystem(new(systems.LogoutSystem))
	}
}

func AddEntityToSector(sectorID core.SectorID) {

}
