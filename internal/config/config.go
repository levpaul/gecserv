package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var (
	GameTickDuration time.Duration
)

func Init() {
	viper.SetConfigName("dev")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	GameTickDuration = time.Second / time.Duration(viper.GetInt("game.tickrate"))
}
