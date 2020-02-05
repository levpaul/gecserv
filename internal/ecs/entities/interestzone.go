package entities

import (
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/ecs/components"
)

type InterestZone struct {
	*BaseEntity
	Entities []core.Entity
	components.Changeable
}
