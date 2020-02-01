package ecs

import "github.com/levpaul/idolscape-backend/internal/ecs/systems"

var (
	pipeErr chan<- error

	sectors []*SectorAdmin
)

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	return nil
}

func initialize() {
	sectors = []*SectorAdmin{
		{id: 1}, // Main and only sector for now
	}

	// Load default systems for all sectors
	for _, s := range sectors {
		s.AddSystem(new(systems.LoginSystem))
		s.AddSystem(new(systems.LogoutSystem))
	}
}

func GetSectors() []*SectorAdmin {
	return sectors
}
