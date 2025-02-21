package main

import (
	"embed"
	"fmt"
	"github.com/je09/spotifind"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"math/rand"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

var configs []spotifind.SpotifindAuth

type config struct {
	SaveLocation string                    `yaml:"saveLocation"`
	Credits      []spotifind.SpotifindAuth `yaml:"credits"`
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "spotifind-gui",
		Width:     400,
		Height:    400,
		MinWidth:  300,
		MinHeight: 350,
		LogLevel:  logger.DEBUG,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func initConfig() (config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("spotifind.yml")
	viper.SetConfigName(".spotifind.yml")

	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/spotifind")
	viper.AddConfigPath("$HOME/.config/spotifind")
	viper.AddConfigPath(".")

	// For Windows
	viper.AddConfigPath(fmt.Sprintf("%s\\AppData\\Roaming\\spotifind", os.Getenv("$HOME")))

	viper.SetEnvPrefix("spotifind")
	viper.BindEnv("spotify_client_id")
	viper.BindEnv("spotify_client_secret")

	if err := viper.ReadInConfig(); err != nil {
		return config{}, err
	}
	var cfg config
	err := viper.Unmarshal(&cfg)
	configs = cfg.Credits

	//randomize configs order
	for i := range configs {
		j := rand.Intn(i + 1)
		configs[i], configs[j] = configs[j], configs[i]

	}

	if err != nil {
		return config{}, err
	}

	return cfg, nil
}
