package systems

import (
	"github.com/levpaul/idolscape-backend/internal/core"
)

type BaseSystem struct {
	sectorID core.SectorID
}

func (bs *BaseSystem) SetSectorID(sectorID core.SectorID) {
	bs.sectorID = sectorID
}

func (bs *BaseSystem) GetSectorID() core.SectorID {
	return bs.sectorID
}
