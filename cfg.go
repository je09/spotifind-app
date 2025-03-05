package main

import (
	"fmt"
	"github.com/je09/spotifind"
	"github.com/spf13/viper"
	"math/rand"
	"os"
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

	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/spotifind")
	viper.AddConfigPath("$HOME/.spotifind")
	viper.AddConfigPath("$HOME/.config/spotifind")
	viper.AddConfigPath(".")

	// For Windows
	viper.AddConfigPath(fmt.Sprintf("%s\\AppData\\Roaming\\spotifind", os.Getenv("$HOME")))

	viper.SetEnvPrefix("spotifind")
	viper.BindEnv("spotify_client_id")
	viper.BindEnv("spotify_client_secret")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	var cfg Config
	err := viper.Unmarshal(&cfg)
	configs = cfg.Credits

	//randomize configs order
	for i := range configs {
		j := rand.Intn(i + 1)
		configs[i], configs[j] = configs[j], configs[i]

	}

	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
