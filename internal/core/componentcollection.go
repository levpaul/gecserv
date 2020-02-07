package core

import (
	"github.com/rs/zerolog/log"
	"reflect"
)

type ComponentCollection []reflect.Type

func NewComponentCollection(components []interface{}) ComponentCollection {
	cc := ComponentCollection{}
	for _, c := range components {
		// Type check the interfaces
		if rawType := reflect.TypeOf(c).Kind(); rawType != reflect.Ptr {
			log.Fatal().
				Str("basic-type", rawType.String()).
				Str("actual-type", reflect.TypeOf(c).String()).
				Msg("Entity system passed an interface which was not a pointer")
		}
		if iType := reflect.TypeOf(c).Elem().Kind(); iType != reflect.Interface {
			log.Fatal().
				Str("basic-type", iType.String()).
				Str("actual-type", reflect.TypeOf(c).Elem().String()).
				Msg("Entity system passed an interface which was not actually an interface!")
		}
		cc = append(cc, reflect.TypeOf(c).Elem())
	}
	return cc
}
