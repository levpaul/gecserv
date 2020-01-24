package proto

import "math/rand"

type CharState struct {
	Color int32
	XPos  float32
	YPos  float32
}

func NewGameState() *CharState {
	return &CharState{
		Color: rand.Int31() % 0xffffff,
		XPos:  rand.Float32() * 20,
		YPos:  rand.Float32() * 20,
	}
}
