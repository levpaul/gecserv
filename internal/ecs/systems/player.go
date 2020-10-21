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

/* PlayerSystem is repsonsible for reading current inputs and moving player
 and updating interest mapping. It also needs to handle out of order messages
 as well as dropped messages and delayed messages.

Basically the PS will need to handle resimulating the world up to THRESHOLD times
per tick.


Simple decision flow can be summarized as follows:
GameTick happens - tick _n_
<LOCK BUS ENQUEUE (bus should stage for next tick)>
1. Pull message from bus - messages should be in a list of priorityQueue based on desc tickSeq per player.
   Messages will contain all actions taken by player between their lastAck and tickSeq
2. If player's last action was simulated; compare actionWindow for mispredictions; IF MISPREDICT:
 - We need to resimulate - the question is how can we minimize what we need to resim here.
 - Worst case - go through all players with simulated actions first; find the lowest tickSeq for mispredictions
     and then simulate the world entirely from that point until n. Then process the rest of the players
 - Better case - only resimulate the player and any other entities player interacted with .... seems like this
     approach will blow up though in dense areas quite quickly. The effort of search here could be overproportional
     compared to simply resimulating
3. ElIf message tickSeq >= n; then simulate action based on worldState n and action at N in the actionWindow,
     stage entity changes
4. ElIf message tickSeq > n; then skip
5. ElIf player has no message - simulate player action based on last input
6. Update player buffer value to send back to player so they can tweak the time dilation on their end
7. Allow interest system to detect changes related to each player for the tick and send updates.
*/
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
