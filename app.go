package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/je09/spotifind-app/common"
	"strings"
	"sync"

	"github.com/je09/spotifind"
	"github.com/labstack/gommon/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	releaseURL            = "https://github.com/je09/spotifind-app/releases/latest"
	installationManualURL = "https://github.com/je09/spotifind-app?tab=readme-ov-file#step-1---installation"
)

// SpotifindApp struct
type SpotifindApp struct {
	ctx context.Context

	h     SpotifindHandler
	cfg   common.ConfigManager
	rls   common.ReleaseManager
	cache Cache

	// searchOnce is used in case of multiple miss-clicks on search button.
	searchOnce sync.Once
}

// NewApp creates a new SpotifindApp application struct
func NewApp() *SpotifindApp {
	return &SpotifindApp{
		cfg:   &common.ConfigManagerImpl{},
		cache: NewCache(),
	}
}

func (a *SpotifindApp) startup(ctx context.Context) {
	a.ctx = ctx

	// Load config and initialize global config for hotswap.
	cfg, cfgs, err := a.cfg.InitConfig()
	if err != nil {
		pb := common.NewPathBuilder()
		a.Alert(fmt.Sprintf("Can't find \"spotifind.yml\" file anywhere. Please check the instruction!\n"+
			"Here are the paths I've tried so far %v", pb.ConfigLocations()))
		runtime.BrowserOpenURL(a.ctx, installationManualURL)
		runtime.LogFatalf(a.ctx, "Error reading config: %v", err)
	}
	configs = cfgs
	runtime.LogInfof(ctx, "Config save location: %s", cfg.SaveLocation)

	// Load previously searched queries and ignores from cache.
	err = a.cache.Load()
	if err != nil {
		runtime.LogFatalf(a.ctx, "Error loading cache: %v", err)
	}

	// New release checker.
	a.rls = common.NewReleaseManager()

	// Initialize spotifind-handler.
	a.h, err = NewHandler(configs[0], cfg.SaveLocation)
	if err != nil {
		runtime.LogFatalf(a.ctx, "Error starting spotifind handler: %v", err)
	}
}

func (a *SpotifindApp) Reconnect() error {
	return a.h.Reconnect(a.ctx)
}

func (a *SpotifindApp) Search(q, ignore, market, csvFileName string) {
	queries := strings.Split(q, ",")
	ignores := strings.Split(ignore, ",")

	if err := a.h.LoadKnownPlaylists(csvFileName); err != nil {
		runtime.LogErrorf(a.ctx, "Error loading playlists CSV file: %v", err)
	}

	runtime.LogInfof(a.ctx, "Searching for: %v, ignoring: %v, on market: %s", queries, ignores, market)
	go a.searchOnce.Do(func() {
		if err := a.h.Search(queries, ignores, market); err != nil {
			runtime.LogInfof(a.ctx, "Error searching on market %s: %v", market, err)
			a.ErrorHandler(err)
		}
	})
}

func (a *SpotifindApp) ProgressBar() {
	res := make(spotifind.ProgressChan)
	go a.h.Progress(res)

	go func() {
		for prog := range res {
			mp, _ := json.Marshal(prog)
			runtime.EventsEmit(a.ctx, "rcv:progress", string(mp))
		}
	}()
}

func (a *SpotifindApp) Results() {
	res := make(spotifind.SpotifindChan)
	go a.h.Results(a.ctx, res)

	go func() {
		for p := range res {
			mp, _ := json.Marshal(p)
			runtime.EventsEmit(a.ctx, "rcv:searchResult", string(mp))
		}
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "Done!",
			Message:       "Search completed!\nAll the results have been saved to the CSV file.\nYou can close the app now.",
			DefaultButton: "OK",
		})
		a.ctx.Done()
	}()
}

func (a *SpotifindApp) Alert(t string) {
	_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Title:   "Oops!",
		Message: t,
	})
	if err != nil {
		log.Errorf("Error showing alert: %v", err)
	}
}

func (a *SpotifindApp) ErrorHandler(err error) {
	if errors.Is(err, spotifind.ErrTimeout) || errors.Is(err, spotifind.ErrTokenExpired) {
		if err = a.Reconnect(); err == nil {
			a.Alert(
				"Spotify said there were too many requests from us.\n" +
					"You can try again later, a few hours later at least.\n" +
					"Or use another auth config.",
			)
		}
		return
	}
	a.Alert(fmt.Sprintf("Something went wrong: %v. "+
		"If the error persists - open an issue here (https://github.com/je09/spotifind-app) and attach the log file (%s)", err, common.LogFileLocation))
}

func (a *SpotifindApp) CheckForNewRelease() {
	release, err := a.rls.NewRelease(Version)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error checking for new release: %v", err)
		return
	}

	if release != "" {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Title:         "New Release Available",
			Message:       fmt.Sprintf("A new release (%s) is available. Please update your application.", release),
			DefaultButton: "Ok",
		})
		runtime.BrowserOpenURL(a.ctx, releaseURL)
	}
}

func (a *SpotifindApp) Markets() []string {
	return spotifind.MarketsAll
}

func (a *SpotifindApp) LoadCachedSearch() []string {
	res := a.cache.PreviousSearch().Searches
	runtime.LogInfof(a.ctx, "Previous searches: %v", res)
	return res
}

func (a *SpotifindApp) LoadCachedIgnore() []string {
	res := a.cache.PreviousSearch().Ignores
	runtime.LogInfof(a.ctx, "Previous ignores: %v", res)
	return res
}

func (a *SpotifindApp) SaveCache(search, ignore string) {
	err := a.cache.Append(search, ignore)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error saving cache: %v", err)
	}
}
