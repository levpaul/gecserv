package systems

import (
	"context"
	"github.com/levpaul/gecserv/internal/config"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/ecs/entities"
	"github.com/levpaul/gecserv/internal/fb"
	"github.com/rs/zerolog/log"
	"math"
)

// TODO: Make the gameConfigs
var (
	moveDecel   float32 = 10
	moveSpeed   float32 = 300
	minMomentum float32 = 0.1
)

// PlayerSystem is repsonsible for reading current inputs and moving player
// and updating interest mapping
type PlayerSystem struct {
	BaseSystem
	events  chan eb.Event
	players map[float64]*entities.PlayerE
}

func (ps *PlayerSystem) Init() {
	ps.events = make(chan eb.Event, 128)
	eb.Subscribe(eb.N_PLAYER_INPUT, ps.events)

	ps.players = ps.sa.GetPlayerListSingleton()
}
func (ps *PlayerSystem) Update(ctx context.Context, dt core.GameTick) {

	// Loop through all moving players and apply player-drag
	ps.applyPlayerDrag(ctx, dt)

	// Loop through player input messages add to player input buffers
	for {
		select {
		case l := <-ps.events:
			switch data := l.Data.(type) {
			case eb.PlayerInputMsg:
				// TODO: Read Seq and interweave appropriately
				ps.handlePlayerInput(data.FromPlayerSID, data.Msg.Actions, data.Msg.CamAngle, dt)
			default:
				log.Error().Interface("type", data).Msg("Unsupported message type received on player channel")
			}
		case <-ctx.Done():
			return
		default:
			return
		}
	}
}

func (ps *PlayerSystem) applyPlayerDrag(ctx context.Context, dt core.GameTick) {
	var tickScaling = float32(dt) / float32(config.TickRate)
	for _, p := range ps.players {
		select {
		case <-ctx.Done():
			return

		default:
			if p.Momentum.X == 0 && p.Momentum.Y == 0 {
				continue
			}

			p.Momentum.X -= p.Momentum.X * moveDecel * tickScaling
			p.Momentum.Y -= p.Momentum.Y * moveDecel * tickScaling

			if float32(math.Abs(float64(p.Momentum.X))) < minMomentum && float32(math.Abs(float64(p.Momentum.Y))) < minMomentum {
				p.Momentum.X = 0
				p.Momentum.Y = 0
			}
		}
	}
}

func (ps *PlayerSystem) handlePlayerInput(sid float64, actions []fb.PlayerAction, camAngle float64, dt core.GameTick) {
	p, ok := ps.players[sid]
	if !ok {
		log.Warn().Float64("sid", sid).Msg("player not found in handlePlayerInput")
		return
	}
	var tickScaling = float32(dt) / float32(config.TickRate)

	moveF := false
	moveB := false
	moveL := false
	moveR := false
	for _, v := range actions {
		switch v {
		case fb.PlayerActionFORWARD:
			moveF = true
		case fb.PlayerActionBACKWARD:
			moveB = true
		case fb.PlayerActionLEFT:
			moveL = true
		case fb.PlayerActionRIGHT:
			moveR = true
		default:
			log.Warn().Str("action", v.String()).Msg("Invalid player action found from player system")
		}
	}

	if !(moveL || moveR || moveF || moveB) {
		return
	}

	dirY := 0
	dirX := 0
	if moveF != moveB {
		if moveF {
			dirY = 1
		} else {
			dirY = -1
		}
	}
	if moveL != moveR {
		if moveL {
			dirX = 1
		} else {
			dirX = -1
		}
	}

	rawForceX := float32(dirX) * moveSpeed * tickScaling
	rawForceY := float32(dirY) * moveSpeed * tickScaling
	cosAngle := float32(math.Cos(camAngle))
	sinAngle := -float32(math.Sin(camAngle))

	p.Momentum.X += rawForceX*cosAngle - rawForceY*sinAngle
	p.Momentum.Y += rawForceY*cosAngle + rawForceX*sinAngle

	p.Position.X += p.Momentum.X * tickScaling
	p.Position.Y += p.Momentum.Y * tickScaling
}
