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
	homeDir, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName("spotifind.yml")

	viper.AddConfigPath(fmt.Sprintf("%s", homeDir))
	viper.AddConfigPath(fmt.Sprintf("%s/spotifind", homeDir))
	viper.AddConfigPath(fmt.Sprintf("%s/Documents/spotifind", homeDir))
	viper.AddConfigPath(fmt.Sprintf("%s/.spotifind", homeDir))
	viper.AddConfigPath(fmt.Sprintf("%s/.config/spotifind", homeDir))
	viper.AddConfigPath(".")

	// For Windows
	viper.AddConfigPath(fmt.Sprintf("%s\\AppData\\Roaming\\spotifind", homeDir))

	viper.SetEnvPrefix("spotifind")
	_ = viper.BindEnv("spotify_client_id")
	_ = viper.BindEnv("spotify_client_secret")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, err
	}
	configs = cfg.Credits

	//randomize configs order
	for i := range configs {
		j := rand.Intn(i + 1)
		configs[i], configs[j] = configs[j], configs[i]

	}

	return cfg, nil
}
