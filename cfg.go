package main

import (
	"github.com/je09/spotifind"
	"github.com/spf13/viper"
	"math/rand"
)

// ConfigManager interface
type ConfigManager interface {
	InitConfig() (Config, error)
}

type Config struct {
	SaveLocation string
	Credits      []spotifind.SpotifindAuth `yaml:"credits"`
}

type ConfigManagerImpl struct{}

func (c *ConfigManagerImpl) InitConfig() (Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("spotifind.yml")

	pb := NewPathBuilder()
	for _, path := range pb.ConfigLocations() {
		viper.AddConfigPath(path)
	}

	viper.SetEnvPrefix("spotifind")
	_ = viper.BindEnv("spotify_client_id")
	_ = viper.BindEnv("spotify_client_secret")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	configs = cfg.Credits

	// Randomize configs order
	for i := range configs {
		j := rand.Intn(i + 1)
		configs[i], configs[j] = configs[j], configs[i]
	}

	return cfg, nil
}
