package entities

import (
	"github.com/levpaul/idolscape-backend/internal/ecs/components"
)

type PlayerE struct {
	*BaseEntity
	components.Position
	components.Changeable
	components.NetworkSession
	components.Color
}
