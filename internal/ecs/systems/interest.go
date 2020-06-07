package systems

import (
	"context"
	"github.com/levpaul/gecserv/internal/core"
	"github.com/levpaul/gecserv/internal/eb"
	"github.com/levpaul/gecserv/internal/ecs/components"
	"github.com/levpaul/gecserv/internal/ecs/entities"
	"github.com/rs/zerolog/log"
)

// InterestSystem is responsible for updating a singleton map of interest buckets
// containing all entities in subdivisions of the sectors map, used by the propagator
// to send relevant map state only to clients. InterestSystem listens for objectMove
// updates from the eventbus and updates all entities from there. May add a scheduled
// full update in too
type InterestSystem struct {
	BaseSystem
	cc     core.ComponentCollection
	events chan eb.Event
}

func (is *InterestSystem) Init() {
	// Create and initialize the interestMap as an entity
	iMEnt := entities.InterestMapE{
		BaseEntity: entities.NewBaseEntity(),
		InterestMap: components.InterestMap{
			Imap: make([][]core.EntityIDs, core.InterestSegmentsX),
		},
	}
	for i := range iMEnt.InterestMap.Imap {
		iMEnt.InterestMap.Imap[i] = make([]core.EntityIDs, core.InterestSegmentsY)
	}
	iMEnt.Lookup = make(map[core.EntityID]core.Vec2Uint8)

	sm := is.sa.GetSectorMap()
	iMEnt.SegSizeX = float32(sm.MaxX) / core.InterestSegmentsX
	iMEnt.SegSizeY = float32(sm.MaxY) / core.InterestSegmentsY

	is.sa.AddEntity(iMEnt)
	is.cc = core.NewComponentCollection([]interface{}{
		new(components.ChangeableComponent),
		new(components.PositionalComponent),
	})

	is.events = make(chan eb.Event, 128)
	eb.Subscribe(eb.S_REMOVED_ENT, is.events)
}

func (is *InterestSystem) Update(ctx context.Context, dt core.GameTick) {
	interestMapEnts := is.sa.FilterEntitiesByCC(core.NewComponentCollection([]interface{}{
		new(components.InterestMapComponent),
	}))
	ime := interestMapEnts.Next()
	if ime == nil || interestMapEnts.Next() != nil {
		log.Fatal().Msg("Unexpected interestmap exception, either interest map not found or more than 1 was!")
	}
	im := ime.(components.InterestMapComponent).GetInterestMap()

	// First check removedEntities queue and update interest map
	for empty := false; !empty; {
		select {
		case rawEvent := <-is.events:
			switch ev := rawEvent.Data.(type) {
			case eb.S_REMOVED_ENT_T:
				handleRemovedEnt(core.EntityID(ev), im)
			default:
				log.Error().Interface("type", rawEvent.Data).Msg("Unsupported message type received on interest channel")
			}
		case <-ctx.Done():
			return
		default:
			empty = true
			break
		}
	}

	// Second loop through all changeable entities w/ position
	// If changed, update interest map w/ new coordinates
	ents := is.sa.FilterEntitiesByCC(is.cc)
	for en := ents.Next(); en != nil; en = ents.Next() {
		// Check if entity changed
		chCp, ok := en.(components.ChangeableComponent)
		if !ok || !chCp.GetChangeable().Changed {
			continue
		}
		chCp.GetChangeable().Changed = false

		// Get relative interest map position
		eid := en.ID()
		posCp, ok := en.(components.PositionalComponent)
		if !ok {
			log.Error().Uint32("entity", uint32(eid)).Msg("Failed to turn entity into position component at interest system")
			continue
		}
		imPosX := uint8(posCp.GetPosition().X / im.SegSizeX)
		imPosY := uint8(posCp.GetPosition().Y / im.SegSizeY)

		// Check to see if entity is new
		old, isInLookup := im.Lookup[eid]
		if !isInLookup {
			im.Imap[imPosX][imPosY] = append(im.Imap[imPosX][imPosY], en.ID())
			im.Lookup[eid] = core.Vec2Uint8{imPosX, imPosY}
			continue
		}

		// Check if no sector position update required
		if old.X == imPosX && old.Y == imPosY {
			continue
		}

		// Update sector in IM by deleting old entry and adding new
		for i, v := range im.Imap[old.X][old.Y] {
			if v == eid {
				im.Imap[old.X][old.Y][i] = im.Imap[old.X][old.Y][len(im.Imap[old.X][old.Y])-1]
				im.Imap[old.X][old.Y] = im.Imap[old.X][old.Y][:len(im.Imap[old.X][old.Y])-1]
				im.Imap[imPosX][imPosY] = append(im.Imap[imPosX][imPosY], eid)
				break
			}
		}
	}
}

func handleRemovedEnt(en core.EntityID, imEn components.InterestMap) {
	pos, ok := imEn.Lookup[en]
	if !ok { // Early exit if entity is not in interest map
		return
	}

	// Delete ent from interest map
	im := imEn.Imap[pos.X][pos.Y]
	for i := range im {
		if im[i] == en {
			im[i] = im[len(im)-1]
			break
		}
	}
	imEn.Imap[pos.X][pos.Y] = im[:len(im)-1]
}
