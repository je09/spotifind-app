package cli

import (
	"fmt"
	"github.com/je09/spotifind"
	"github.com/je09/spotifind-app/common"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spotyfind2",
	Short: "Spotifind2 is a CLI tool for searching Spotify playlists",
}
var spotifindHandler *SpotifyHandler
var configs []spotifind.SpotifindAuth

type credsType struct {
	Credits []spotifind.SpotifindAuth `yaml:"credits"`
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cfg := common.ConfigManagerImpl{}
	_, cfgs, err := cfg.InitConfig()
	configs = cfgs

	if err != nil {
		fmt.Println(err)
		return
	}
	spotifindHandler, err = NewSpotifyHandler()
	if err != nil {
		panic(err)
	}
}
