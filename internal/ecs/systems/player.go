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

func (ps *PlayerSystem) handlePlayerInput(sid float64, actions []fb.PlayerAction, camAngle float64, dt core.GameTick) {
	p, ok := ps.players[sid]
	if !ok {
		log.Warn().Float64("sid", sid).Msg("player not found in handlePlayerInput")
		return
	}
	if p.Momentum.X == 0 && p.Momentum.Y == 0 && len(actions) == 0 {
		return
	}
	if p.Momentum.X < 0.0001 && p.Momentum.Y < 0.0001 && len(actions) == 0 {
		p.Momentum.X = 0
		p.Momentum.Y = 0
		return
	}

	var moveDecel float32 = 10
	var moveSpeed float32 = 300
	var tickScaling = float32(dt) / float32(config.TickRate)

	p.Momentum.X -= p.Momentum.X * moveDecel * tickScaling
	p.Momentum.Y -= p.Momentum.Y * moveDecel * tickScaling

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

	// DEBUG

	newForceX := float32(dirX) * moveSpeed * tickScaling
	newForceY := float32(dirY) * moveSpeed * tickScaling
	cosAngle := float32(math.Cos(camAngle))
	sinAngle := float32(math.Sin(camAngle))
	newForceX = newForceX*cosAngle - newForceY*sinAngle
	newForceY = newForceY*cosAngle + newForceX*sinAngle

	log.Info().Float64("Angle", camAngle).Float32("xforce", newForceX).Float32("yforce", newForceY).
		Float32("xmomen", p.Momentum.X).Float32("ymomen", p.Momentum.Y).Send()

	p.Momentum.X += newForceX
	p.Momentum.Y += newForceY

	p.Position.X += p.Momentum.X * tickScaling
	p.Position.Y += p.Momentum.Y * tickScaling

	// TODO: Allow camera rotation relative WASD controls
	//this.vec.setFromMatrixColumn( this.oC.object.matrix, 0 );
	//this.game.char.position.addScaledVector( this.vec, this.velocity.x * timeDelta);
	//this.vec.crossVectors( this.oC.object.up, this.vec );
	//this.game.char.position.addScaledVector( this.vec, this.velocity.z * timeDelta);
}
