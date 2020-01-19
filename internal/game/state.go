package game

import "math/rand"

type GameState struct {
	charX float32
	charY float32
}

func NewGameState() *GameState {
	return &GameState{
		charX: rand.Float32() * 200,
		charY: rand.Float32() * 200,
	}
}

func (gs *GameState) GetPos() (float32, float32) {
	return gs.charX, gs.charY
}
