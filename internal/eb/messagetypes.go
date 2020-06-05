package eb

import (
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/ecs/entities"
)

type PlayerSyncMessage struct {
	ToPlayerSID float64
	Players     []*entities.PlayerE
	Tick        core.GameTick
}
