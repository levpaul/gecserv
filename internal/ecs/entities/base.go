package entities

import (
	"github.com/levpaul/idolscape-backend/internal/core"
	"sync"
)

var (
	idCounter core.EntityID
	counterMu sync.Mutex
)

type BaseEntity struct {
	id core.EntityID
	n  core.Entity
	p  core.Entity
}

func (e *BaseEntity) ID() core.EntityID {
	return e.id
}

func (e *BaseEntity) Next() core.Entity {
	return e.n
}

func (e *BaseEntity) SetNext(next core.Entity) {
	e.n = next
}

func (e *BaseEntity) Prev() core.Entity {
	return e.p
}

func (e *BaseEntity) SetPrev(prev core.Entity) {
	e.p = prev
}

func NewBaseEntity() *BaseEntity {
	be := new(BaseEntity)
	be.id = newEntityID()
	return be
}

func newEntityID() core.EntityID {
	counterMu.Lock()
	defer counterMu.Unlock()
	idCounter += 1
	return idCounter
}
