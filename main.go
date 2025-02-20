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

type credsType struct {
	Credits []spotifind.SpotifindAuth `yaml:"credits"`
}

func main() {
	initConfig()

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

func initConfig() {
	c := "spotifind.yml"
	viper.SetConfigType("yaml")
	viper.SetConfigName(c)

	//home, _ := os.UserHomeDir()
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("spotifind")
	viper.BindEnv("spotify_client_id")
	viper.BindEnv("spotify_client_secret")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading configs file:", err)
		os.Exit(1)
	}
	var creds credsType
	err := viper.Unmarshal(&creds)
	configs = creds.Credits

	// randomize configs order
	for i := range configs {
		j := rand.Intn(i + 1)
		configs[i], configs[j] = configs[j], configs[i]

	}

	if err != nil {
		return
	}
}
