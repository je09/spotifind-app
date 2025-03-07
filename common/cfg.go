package common

import (
	"github.com/je09/spotifind"
	"github.com/spf13/viper"
	"math/rand"
)

// ConfigManager interface
type ConfigManager interface {
	InitConfig() (Config, []spotifind.SpotifindAuth, error)
}

type Config struct {
	SaveLocation string
	Credits      []spotifind.SpotifindAuth `yaml:"credits"`
}

type ConfigManagerImpl struct{}

func (c *ConfigManagerImpl) InitConfig() (Config, []spotifind.SpotifindAuth, error) {
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
		return Config{}, nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, nil, err
	}

	cfgs := cfg.Credits

	// Randomize configs order
	for i := range cfgs {
		j := rand.Intn(i + 1)
		cfgs[i], cfgs[j] = cfgs[j], cfgs[i]
	}

	return cfg, cfgs, nil
}
