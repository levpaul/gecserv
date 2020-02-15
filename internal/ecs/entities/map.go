package entities

import (
	"github.com/levpaul/gecserv/internal/ecs/components"
)

type MapE struct {
	*BaseEntity
	components.MapComponent
}

func NewDefaultMap() *MapE {
	me := &MapE{
		BaseEntity: NewBaseEntity(),
		MapComponent: &components.Map{
			Name:    "default map",
			GlobalX: 0,
			GlobalY: 0,
			MaxX:    1000,
			MaxY:    1000,
		},
	}
	return me
}
