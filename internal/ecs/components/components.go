package components

import (
	"github.com/levpaul/idolscape-backend/internal/core"
	"github.com/levpaul/idolscape-backend/internal/fb"
)

// ==================================================================
type Position fb.Vec2T

func (p *Position) GetPosition() *Position { return p }

type PositionComponent interface {
	GetPosition() *Position
}

// ==================================================================
type NetworkSession struct {
	Sid float64
}

type NetworkSessionComponent interface {
	GetNetorkSession() *NetworkSession
}

func (n *NetworkSession) GetNetworkSession() *NetworkSession { return n }

// ==================================================================
type Color struct {
	Col fb.Color
}

type ColorComponent interface {
	GetColor() *Color
}

func (c *Color) GetColor() *Color {
	return c
}

// ==================================================================
type Changeable struct {
	Changed bool
}

type ChangeableComponent interface {
	GetChangeable() *Changeable
}

func (c *Changeable) GetChangeable() *Changeable { return c }

// ==================================================================
type StateHistory struct {
	LastAck core.GameTick
}

type StateHistoryComponent interface {
	GetStateHistory() *StateHistory
}

func (s *StateHistory) GetStateHistory() *StateHistory { return s }

// ==================================================================
