package entities

import (
	"github.com/levpaul/gecserv/internal/ecs/components"
)

type PlayerE struct {
	*BaseEntity
	components.Position
	components.Changeable
	components.NetworkSession
	components.StateHistory
	components.Color
}
