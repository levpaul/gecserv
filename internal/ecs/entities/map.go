package entities

import (
	"github.com/levpaul/idolscape-backend/internal/fb"
)

type MapE struct {
	*BaseEntity
	*fb.MapT
}

func NewDefaultMap() *MapE {
	me := &MapE{
		BaseEntity: NewBaseEntity(),
		MapT: &fb.MapT{
			Name:    "default map",
			GlobalX: 0,
			GlobalY: 0,
			MaxX:    100,
			MaxY:    100,
		},
	}
	return me
}
