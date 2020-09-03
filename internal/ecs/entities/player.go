package entities

import (
	"github.com/levpaul/gecserv/internal/ecs/components"
	"github.com/levpaul/gecserv/internal/fb"
)

type PlayerE struct {
	*BaseEntity
	components.Position
	components.Momentum
	components.Changeable
	components.NetworkedSession
	components.StateHistory
	components.Colored
}

// ToPublicFB returns a generic representation of a Player
// suitable for message transfer
func (p *PlayerE) ToPublicFB() *fb.PlayerT {
	return &fb.PlayerT{
		Posx: p.Position.X,
		Posy: p.Position.Y,
		Momx: p.Momentum.X,
		Momy: p.Momentum.Y,
		Sid:  p.Sid,
		Col:  p.Col,
	}
}
