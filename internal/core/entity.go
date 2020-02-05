package core

type Entity interface {
	ID() EntityID
	Next() Entity
	Prev() Entity
	SetNext(Entity)
	SetPrev(Entity)
}
