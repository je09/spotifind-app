package main

import (
	"embed"
	"fmt"
	"github.com/je09/spotifind"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var Version = "v0.0.1"

//go:embed all:frontend/dist
var assets embed.FS

// for hotswap.
var configs []spotifind.SpotifindAuth

func main() {
	l := NewLogger()
	l.Info(fmt.Sprintf("Starting spotifind-gui version: %s", Version))

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:              "spotifind-gui",
		Width:              400,
		Height:             400,
		MinWidth:           300,
		MinHeight:          350,
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.INFO,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Logger: l,
	})
	if err != nil {
		l.Fatal(fmt.Sprintf("Failed to start spotifind: %v", err))
	}
}
