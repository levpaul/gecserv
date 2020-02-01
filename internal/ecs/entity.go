package ecs

type EntityID uint32

type Entity struct {
	ID         EntityID
	Components []Component
}
