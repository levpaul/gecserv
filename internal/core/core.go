package core

import uuid2 "github.com/google/uuid"

type GameTick uint32
type SectorID uint16
type EntityID uint32

type SenderCloser interface {
	Send([]byte) error
	Close() error
}

type AvatarPubConn struct {
	AID  uuid2.UUID // Avatar ID
	Conn SenderCloser
}
