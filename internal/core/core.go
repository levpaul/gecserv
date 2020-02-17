package core

import (
	"fmt"
	"math"
)

type GameTick uint32
type SectorID uint16

type SenderCloser interface {
	Send([]byte) error
	Close() error
}

type SessionPubConn struct {
	SID  float64
	Conn SenderCloser
}

func SIDStr(sid float64) string {
	return fmt.Sprintf("%x", math.Float64bits(sid))
}

type Vec2Uint8 struct {
	X uint8
	Y uint8
}
