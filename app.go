package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/je09/spotifind"
	"github.com/je09/spotifind-app/pkg/csv"
	"github.com/labstack/gommon/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SpotifindService interface
type SpotifindService interface {
	SearchPlaylistPopular(chan spotifind.SpotifindChan, chan spotifind.ProgressChan, []string, []string) error
	SearchPlaylistUnpopular(chan spotifind.SpotifindChan, chan spotifind.ProgressChan, []string, []string) error
	SearchPlaylistForMarket(chan spotifind.SpotifindChan, chan spotifind.ProgressChan, string, []string, []string) error
	Reconnect(Config) error
	Continue(chan spotifind.SpotifindChan, chan spotifind.ProgressChan) error
}

// SpotifindApp struct
type SpotifindApp struct {
	s             spotifind.Spotifinder
	csv           *csv.CsvHandler
	configManager ConfigManager

	KnownPlaylists []string
	currentConfig  int

	playlistChan spotifind.SpotifindChan
	progChan     spotifind.ProgressChan

	searchOnce sync.Once
	configOnce sync.Once

	ctx context.Context
}

// NewApp creates a new SpotifindApp application struct
func NewApp() *SpotifindApp {
	return &SpotifindApp{
		configManager: &ConfigManagerImpl{},
		playlistChan:  make(spotifind.SpotifindChan),
		progChan:      make(spotifind.ProgressChan),
	}
}

func (a *SpotifindApp) startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := a.configManager.InitConfig()
	if err != nil {
		a.Alert(fmt.Sprintf("Error reading config: %v", err))
		panic(err)
	}
	a.csv = csv.NewCsvHandler(cfg.SaveLocation)

	a.s, err = spotifind.NewSpotifind(configs[a.currentConfig], false)
	if err != nil {
		panic(err)
	}
}

func (a *SpotifindApp) Reconnect() error {
	for {
		if a.currentConfig+1 >= len(configs) {
			return fmt.Errorf("no more configs to try")
		}

		log.Infof("Reconnecting to Spotify with new config: %s", configs[a.currentConfig].ClientID)

		a.currentConfig++
		err := a.s.Reconnect(configs[a.currentConfig])
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error reconnecting to Spotify: %v", err)
		}
		err = a.s.Continue(a.playlistChan, a.progChan)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error continuing Spotify: %v", err)
		}
	}
}

func (a *SpotifindApp) Search(q, ignore, market, csvFileName string) string {
	queries := strings.Split(q, ",")
	ignores := strings.Split(ignore, ",")

	a.csv.Path = fmt.Sprintf("%s/%s.csv", a.csv.Path, csvFileName)
	var err error
	a.KnownPlaylists, err = a.csv.ReadFromFile()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error reading CSV file: %v", err)
	}

	runtime.LogInfof(a.ctx, "Searching for: %v, ignoring: %v, on market: %s", queries, ignores, market)

	go a.searchOnce.Do(func() {
		switch market {
		case "popular":
			err := a.s.SearchPlaylistPopular(a.playlistChan, a.progChan, queries, ignores)
			if err != nil {
				a.ErrorHandler(err)
				log.Errorf("Error searching popular playlists: %v", err)
			}
		case "unpopular":
			err := a.s.SearchPlaylistUnpopular(a.playlistChan, a.progChan, queries, ignores)
			if err != nil {
				a.ErrorHandler(err)
				log.Errorf("Error searching unpopular playlists: %v", err)
			}
		default:
			err := a.s.SearchPlaylistForMarket(a.playlistChan, a.progChan, market, queries, ignores)
			if err != nil {
				a.ErrorHandler(err)
				log.Errorf("Error searching playlists for market %s: %v", market, err)
			}
		}
	})

	return fmt.Sprintf("Query: %s, market: %s, ignore: %s", q, market, ignore)
}

func (a *SpotifindApp) ReturnProgress() {
	go func() {
		for prog := range a.progChan {
			mp, _ := json.Marshal(prog)
			runtime.EventsEmit(a.ctx, "rcv:progress", string(mp))
		}
	}()
}

func (a *SpotifindApp) ReturnResults() {
	go func() {
		for p := range a.playlistChan {
			if a.IsPlaylistKnown(p.Playlist.ExternalURLs["spotify"]) {
				runtime.LogInfof(a.ctx, "Playlist already known: %s", p.Playlist.ExternalURLs["spotify"])
				continue
			}
			if p.Playlist.Contacts == nil {
				p.Playlist.Contacts = []string{"N/A"}
			}
			a.csv.WriteToFile(p.Playlist)
			mp, _ := json.Marshal(p)
			runtime.EventsEmit(a.ctx, "rcv:searchResult", string(mp))
		}
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "Done!",
			Message:       "Search completed!\nAll the results have been saved to the CSV file.\nYou can close the app now.",
			DefaultButton: "OK",
		})
		a.ctx.Done()
	}()
}

func (a *SpotifindApp) Markets() []string {
	return spotifind.MarketsAll
}

func (a *SpotifindApp) Alert(t string) {
	_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Ooops!",
		Message: t,
	})
	if err != nil {
		log.Errorf("Error showing alert: %v", err)
	}
}

func (a *SpotifindApp) ErrorHandler(err error) {
	if err == spotifind.ErrTimeout || err == spotifind.ErrTokenExpired {
		if err = a.Reconnect(); err == nil {
			a.Alert(
				"Spotify said there were too many requests from us.\n" +
					"You can try again later, a few hours later at least.\n" +
					"Or use another auth config.",
			)
		}
		return
	}
	a.Alert(err.Error())
}

func (s *SpotifindApp) IsPlaylistKnown(externalURL string) bool {
	if len(s.KnownPlaylists) == 0 {
		return false
	}

	for _, p := range s.KnownPlaylists {
		if p == externalURL {
			return true
		}
	}
	return false
}
