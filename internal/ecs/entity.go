package ecs

type EntityID uint32

type Entity interface {
	ID() EntityID
}
