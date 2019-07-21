package config

import (
	"github.com/spf13/viper"
)

// Load tells viper to load from the environment
func Load() *viper.Viper {
	viper.AutomaticEnv()
	return viper.GetViper()
}
