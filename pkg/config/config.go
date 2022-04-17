package config

import (
	"github.com/spf13/viper"
)

// Load tells viper to load from the environment
func Load() *viper.Viper {
	SetDefaults()
	viper.AutomaticEnv()

	return viper.GetViper()
}
