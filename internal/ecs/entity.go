package ecs

import (
	"github.com/levpaul/idolscape-backend/internal/core"
)

type Entity interface {
	ID() core.EntityID
}
