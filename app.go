package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/je09/spotifind"
	"github.com/je09/spotifind-app/pkg/csv"
	"github.com/labstack/gommon/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
	"sync"
)

// SpotifindApp struct
type SpotifindApp struct {
	s   spotifind.Search
	csv csv.CsvHandler

	KnownPlaylists []string
	currentConfig  int

	playlistChan spotifind.SpotifindChan
	progChan     spotifind.ProgressChan

	once sync.Once

	ctx context.Context
}

// NewApp creates a new SpotifindApp application struct
func NewApp() *SpotifindApp {
	return &SpotifindApp{
		playlistChan: make(spotifind.SpotifindChan),
		progChan:     make(spotifind.ProgressChan),
	}
}

func (a *SpotifindApp) startup(ctx context.Context) {
	var err error
	a.ctx = ctx
	a.s, err = spotifind.NewSpotifind(configs[a.currentConfig], false)
	if err != nil {
		panic(err)
	}
}

func (a *SpotifindApp) Reconnect() error {
	if a.currentConfig+1 >= len(configs) {
		return fmt.Errorf("no more configs to try")
	}

	log.Infof("Reconnecting to Spotify with new config: %s", configs[a.currentConfig].ClientID)

	a.currentConfig++
	spotifind, err := spotifind.NewSpotifind(configs[a.currentConfig], false)
	if err != nil {
		return err
	}
	a.s = spotifind
	return nil
}

func (a *SpotifindApp) Search(q, ignore, market string) string {
	queries := strings.Split(q, ",")
	ignores := strings.Split(ignore, ",")

	runtime.LogInfof(a.ctx, "Searching for: %v, ignoring: %v, on market: %s", queries, ignores, market)

	a.once.Do(func() {
		switch market {
		case "popular":
			go func() {
				err := a.s.SearchPlaylistPopular(a.playlistChan, a.progChan, queries, ignores)
				a.Alert(err.Error())
				if err != nil {
					log.Errorf("Error searching popular playlists: %v", err)
				}
			}()
		case "unpopular":
			go func() {
				err := a.s.SearchPlaylistUnpopular(a.playlistChan, a.progChan, queries, ignores)
				a.Alert(err.Error())
				if err != nil {
					log.Errorf("Error searching unpopular playlists: %v", err)
				}
			}()
		default:
			go func() {
				err := a.s.SearchPlaylistForMarket(a.playlistChan, a.progChan, market, queries, ignores)
				a.Alert(err.Error())
				if err != nil {
					log.Errorf("Error searching playlists for market %s: %v", market, err)
				}
			}()
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
			mp, _ := json.Marshal(p)
			runtime.EventsEmit(a.ctx, "rcv:searchResult", string(mp))
		}
	}()
}

func (a *SpotifindApp) Markets() []string {
	return spotifind.MarketsAll
}

func (a *SpotifindApp) Alert(t string) {
	_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Alert!",
		Message: t,
	})
	if err != nil {
		log.Errorf("Error showing alert: %v", err)
	}
}

func (a *SpotifindApp) WindowSize(h, w int) {
	runtime.WindowSetSize(a.ctx, w, h)
}
