package entities

import (
	"github.com/levpaul/idolscape-backend/internal/fb"
)

type PlayerE struct {
	*BaseEntity
	*fb.PlayerT
}
