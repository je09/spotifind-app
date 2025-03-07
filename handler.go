package main

import (
	"context"
	"errors"
	"github.com/je09/spotifind"
	"github.com/je09/spotifind-app/pkg/csv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	errNoMoreConfigs = errors.New("no more configs to try")
)

const (
	Popular   = "popular"
	Unpopular = "unpopular"
)

type SpotifindHandler interface {
	Search(q, i []string, m string) error
	Reconnect(ctx context.Context) error
	LoadKnownPlaylists(fileName string) error
	Results(ctx context.Context, p spotifind.SpotifindChan)
	Progress(p spotifind.ProgressChan)
}

type Handler struct {
	s              spotifind.Spotifinder
	csv            csv.ReaderWriter
	knownPlaylists map[string]struct{}
	currCfg        int

	plChan spotifind.SpotifindChan
	prChan spotifind.ProgressChan
}

func NewHandler(auth spotifind.SpotifindAuth, resSavePath string) (SpotifindHandler, error) {
	s, err := spotifind.NewSpotifind(auth, false)
	if err != nil {
		return nil, err
	}

	return &Handler{
		s:              s,
		csv:            csv.New(resSavePath),
		plChan:         make(spotifind.SpotifindChan),
		prChan:         make(spotifind.ProgressChan),
		knownPlaylists: make(map[string]struct{}),
	}, nil
}

func (h *Handler) Search(q, i []string, m string) error {
	switch m {
	case Popular:
		err := h.s.SearchPlaylistPopular(h.plChan, h.prChan, q, i)
		if err != nil {
			return err
		}
	case Unpopular:
		err := h.s.SearchPlaylistUnpopular(h.plChan, h.prChan, q, i)
		if err != nil {
			return err
		}
	default: // specific market,
		err := h.s.SearchPlaylistForMarket(h.plChan, h.prChan, m, q, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) Results(ctx context.Context, p spotifind.SpotifindChan) {
	for pl := range h.plChan {
		if h.isKnown(pl.Playlist.ExternalURLs["spotify"]) {
			runtime.LogInfof(ctx, "Playlist already known: %s", pl.Playlist.ExternalURLs["spotify"])
			continue
		}

		if pl.Playlist.Contacts == nil {
			pl.Playlist.Contacts = []string{"N/A"}
		}

		if err := h.csv.WriteToFile(pl.Playlist); err != nil {
			runtime.LogErrorf(ctx, "Error writing playlist to file: %v", err)
		}
		p <- pl
	}
	close(p)
}

func (h *Handler) Progress(p spotifind.ProgressChan) {
	for prog := range h.prChan {
		p <- prog
	}
	close(p)
}

func (h *Handler) LoadKnownPlaylists(fileName string) error {
	if err := h.csv.SetFilePath(fileName); err != nil {
		return err
	}

	var err error
	h.knownPlaylists, err = h.csv.ReadFromFile()
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) Reconnect(ctx context.Context) error {
	for {
		if h.currCfg+1 >= len(configs) {
			return errNoMoreConfigs
		}

		runtime.LogInfof(ctx, "Reconnecting to Spotify with new config")

		h.currCfg++
		err := h.s.Reconnect(configs[h.currCfg])
		if err != nil {
			runtime.LogErrorf(ctx, "Error reconnecting to Spotify: %v", err)
		}
		err = h.s.Continue(h.plChan, h.prChan)
		if err != nil {
			runtime.LogErrorf(ctx, "Error continuing Spotify: %v", err)
		}
	}
}

func (h *Handler) isKnown(playlistURI string) bool {
	if _, ok := h.knownPlaylists[playlistURI]; ok {
		return true
	}

	return false
}
