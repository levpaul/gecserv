package ecs

import (
	"github.com/levpaul/idolscape-backend/internal/core"
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

//// TODO: This may be a bit loose - we could instead have a reverse lookup
////   of Systems -> Sectors
//func AddEntityToSector(en Entity, sectorID core.SectorID) error {
//	sa, ok := sectors[sectorID]
//	if !ok {
//		return errors.New("invalid sector ID")
//	}
//
//	sa.addEntity(en)
//	return nil
//}

//
//func AddSystemToSector(sys System, sectorID core.SectorID) error {
//	sa, ok := sectors[sectorID]
//	if !ok {
//		return errors.New("invalid sector ID")
//	}
//
//	sa.addSystem(sys)
//	return nil
//}
//
//func AddEntitySystemToSector(sys EntitySystem, ifce interface{}, sectorID core.SectorID) error {
//	sa, ok := sectors[sectorID]
//	if !ok {
//		return errors.New("invalid sector ID")
//	}
//
//	sa.addEntitySystem(sys, ifce)
//	return nil
//}

func AddNewSector() *SectorAdmin {
	sa := newSectorAdmin()
	sectors[sa.id] = sa
	return sa
}

//func RemoveEntityFromSector(en core.EntityID, sectorID core.SectorID) {
//	sectors[sectorID].removeEntity(en)
//}
//
//func GetEntityFromSector(entityID core.EntityID, sectorID core.SectorID) (Entity, error) {
//	sa, ok := sectors[sectorID]
//	if !ok {
//		return nil, errors.New("invalid sector ID")
//	}
//
//	ent := sa.getEntity(entityID)
//	if ent == nil {
//		return nil, errors.New("entity not found in sector")
//	}
//	return ent, nil
//}
