package core

type EntityID uint32
type EntityIDs []EntityID

type Entity interface {
	ID() EntityID
}
