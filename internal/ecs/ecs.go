package ecs

import (
	"github.com/levpaul/gecserv/internal/core"
)

var (
	pipeErr chan<- error

	sectors map[core.SectorID]*SectorAdmin
)

func Start(pErr chan<- error) error {
	pipeErr = pErr
	initialize()
	go updateLoop()
	return nil
}

func initialize() {
	sectors = make(map[core.SectorID]*SectorAdmin)
}

func AddNewSector() *SectorAdmin {
	sa := newSectorAdmin()
	sectors[sa.id] = sa
	return sa
}
