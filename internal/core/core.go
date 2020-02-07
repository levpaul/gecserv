package core

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"math"
)

type GameTick uint32
type SectorID uint16

type SenderCloser interface {
	Send([]byte) error
	Close() error
}

type AvatarPubConn struct {
	AID  uuid2.UUID // Avatar ID
	Conn SenderCloser
}

func SIDStr(sid float64) string {
	return fmt.Sprintf("%x", math.Float64bits(sid))
}
