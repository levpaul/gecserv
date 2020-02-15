package components

import (
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/fb"
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
	GetNetworkSession() *NetworkSession
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
type InterestMap struct {
	Imap     [][]core.EntityIDs
	Lookup   map[core.EntityID]core.Vec2Uint8
	SegSizeX float32
	SegSizeY float32
}

type InterestMapComponent interface {
	GetInterestMap() InterestMap
}

func (im InterestMap) GetInterestMap() InterestMap { return im }

func (im InterestMap) GetPosIMCoords(p fb.Vec2T) core.Vec2Uint8 {
	return core.Vec2Uint8{
		X: uint8(p.X / im.SegSizeX),
		Y: uint8(p.Y / im.SegSizeY),
	}
}

// ==================================================================
type Map fb.MapT

type MapComponent interface {
	GetMap() *Map
}

func (m *Map) GetMap() *Map { return m }

// ==================================================================
