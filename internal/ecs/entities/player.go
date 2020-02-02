package entities

import (
	"github.com/levpaul/idolscape-backend/internal/ecs"
	"github.com/levpaul/idolscape-backend/internal/fb"
)

type PlayerE struct {
	*ecs.BaseEntity
	*fb.PlayerT
}
