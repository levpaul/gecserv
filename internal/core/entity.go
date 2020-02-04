package core

type Entity interface {
	ID() EntityID
	Next() Entity
	SetNext(Entity)
	Prev() Entity
	SetPrev(Entity)
}
