package entities

import "github.com/levpaul/idolscape-backend/internal/ecs/components"

type InterestMapE struct {
	*BaseEntity
	components.InterestMap
}
