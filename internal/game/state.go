package game

import "math/rand"

type CharState struct {
	color int32
	xPos  float32
	yPos  float32
}

func NewGameState() *CharState {
	return &CharState{
		color: rand.Int31() % 0xffffff,
		xPos:  rand.Float32() * 20,
		yPos:  rand.Float32() * 20,
	}
}

func (gs *CharState) GetPos() (float32, float32) {
	return gs.xPos, gs.yPos
}

func (gs *CharState) GetCol() int32 {
	return gs.color
}
