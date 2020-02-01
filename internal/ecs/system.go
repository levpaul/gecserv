package ecs

import (
	"context"
	"github.com/levpaul/idolscape-backend/internal/core"
)

type System interface {
	Update(ctx context.Context, delta core.GameTick)
}

type Initializer interface {
	Init()
}
