package ecs

import (
	"errors"
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
		// TODO: Split these out to some sort of loader and make AddSystem/AddSector public
		s.addSystem(new(systems.LoginSystem))
		s.addSystem(new(systems.LogoutSystem))
	}
}

// TODO: This may be a bit loose - we could instead have a reverse lookup
//   of Systems -> Sectors
func AddEntityToSector(sectorID core.SectorID, en Entity) error {
	sa, ok := sectors[sectorID]
	if !ok {
		return errors.New("invalid sector ID")
	}

	sa.addEntity(en)
	return nil
}
