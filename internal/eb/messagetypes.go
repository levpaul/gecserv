package eb

import (
	uuid2 "github.com/google/uuid"
	"github.com/levpaul/gecserv/internal/fb"
)

type PlayerSyncMessage struct {
	ToPlayer uuid2.UUID
	Players  []*fb.PlayerT
}
