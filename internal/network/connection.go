package network

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/idolscape-backend/internal/game"
	"github.com/levpaul/idolscape-backend/internal/state"
	"github.com/pion/webrtc"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
)

var conns []Connection
var connsLock sync.Mutex

var newConns []*Connection
var newConnLock sync.Mutex
var discConns []*Connection
var discConnLock sync.Mutex

type Connection struct {
	Uuid uuid2.UUID
	pc   *webrtc.PeerConnection
	dc   *webrtc.DataChannel
	GS   *game.CharState
}

func NewConnection(pc *webrtc.PeerConnection, dc *webrtc.DataChannel) *Connection {
	connsLock.Lock()
	defer connsLock.Unlock()
	conns = append(conns, Connection{
		Uuid: uuid2.New(),
		pc:   pc,
		dc:   dc,
		GS:   game.NewGameState(),
	})

	newConnLock.Lock()
	defer newConnLock.Unlock()
	newConns = append(newConns, &conns[len(conns)-1])

	return &conns[len(conns)-1]
}

func (c *Connection) SendNewPlayer(p *Connection) {
	c.dc.SendText(fmt.Sprintf(`{"type": "newchar", "color": %d, "x":%f,"y":%f, "uuid": "%s"}`, p.GS.Color, p.GS.XPos, p.GS.YPos, p.Uuid))
}

func (c *Connection) SendDisconnectedPlayer(p *Connection) {
	c.dc.SendText(fmt.Sprintf(`{"type": "discchar", "uuid": "%s"}`, p.Uuid))
}

func (c *Connection) Disconnect() {
	log.Info().Str("UUID", c.Uuid.String()).Msg("Disconnecting client")
	connsLock.Lock()
	defer connsLock.Unlock()

	// Delete conns
	for i := range conns {
		if conns[i].Uuid == c.Uuid {
			conns[i] = conns[len(conns)-1]
			conns = conns[:len(conns)-1]
			break
		}
	}

	discConnLock.Lock()
	defer discConnLock.Unlock()
	discConns = append(discConns, c)
}

func (c *Connection) SendInitState() {
	c.dc.SendText(fmt.Sprintf(`{"type": "initpos", "color": %d, "x":%f,"y":%f, "uuid": "%s"}`, c.GS.Color, c.GS.XPos, c.GS.YPos, c.Uuid))
	c.SendOtherCharsState()
}

func (c *Connection) SendOtherCharsState() {
	// TODO: can make this more efficient O(1)?
	var otherConns []*Connection
	connsLock.Lock()
	for i := range conns {
		oc := &conns[i]
		if c.Uuid.String() != oc.Uuid.String() {
			otherConns = append(otherConns, oc)
		}
	}
	connsLock.Unlock()

	b, err := json.Marshal(struct {
		Type  string
		Chars []*Connection
	}{
		Type:  "charlist",
		Chars: otherConns,
	})
	if err != nil {
		log.Err(err).Msg("Error marshalling channel list to send back")
		return
	}

	// TODO: Replace all this stuff with protobuffs or something
	c.dc.SendText(strings.ToLower(string(b)))
}

// ===========================================================================
// Protobuff methods

func (c *Connection) SendNewPlayerState(s *state.State) {
	b, err := proto.Marshal(s)
	if err != nil {
		log.Err(err).Msg("Failed to marshal proto for new player state")
	}
	c.dc.Send(b)
}

// ===========================================================================

type MoveMessage struct {
	X float32
	Y float32
}
