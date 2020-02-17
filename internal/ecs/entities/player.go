package entities

import (
	"github.com/levpaul/gecserv/internal/ecs/components"
	"github.com/levpaul/gecserv/internal/fb"
)

type PlayerE struct {
	*BaseEntity
	components.Position
	components.Changeable
	components.NetworkSession
	components.StateHistory
	components.Color
}

func (p *PlayerE) ToFB() *fb.PlayerT {
	return &fb.PlayerT{
		Posx: p.X,
		Posy: p.Y,
		Sid:  p.Sid,
		Col:  p.Col,
	}
}
