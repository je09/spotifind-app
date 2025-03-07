package cli

import (
	"errors"
	"fmt"
	"github.com/je09/spotifind"
	"github.com/je09/spotifind-app/pkg/csv"
	"math"
	"strings"
	"time"
)

const (
	printFormat = "%s - 🙏(%s) ❤️%d 🎸%s (%s)\n" //playlistName, contacts, followers, styles, region
)

type SpotifyHandler struct {
	spotifind spotifind.Search
	Csv       csv.CSV

	KnownPlaylists map[string]struct{}
	currentConfig  int
}

func NewSpotifyHandler() (*SpotifyHandler, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("no config found")
	}

	rootCmd.Printf("\r"+Green+"using conf: %s\n"+Reset, configs[0].ClientID)
	spotifind, err := spotifind.NewSpotifind(configs[0], false)
	if err != nil {
		return nil, err
	}

	return &SpotifyHandler{
		spotifind: spotifind,
		Csv:       csv.CSV{},
	}, nil
}

func (s *SpotifyHandler) Reconnect() error {
	if s.currentConfig+1 >= len(configs) {
		return fmt.Errorf("no more configs to try")
	}

	s.currentConfig++
	spotifind, err := spotifind.NewSpotifind(configs[s.currentConfig], false)
	if err != nil {
		return err
	}
	rootCmd.Printf("\r"+Green+"setting new conf: %s\n"+Reset, configs[s.currentConfig].ClientID)
	s.spotifind = spotifind
	return nil
}

func (s *SpotifyHandler) PrintFormattedPlaylist(playlist spotifind.Playlist) {
	rootCmd.Printf("\r")
	rootCmd.Printf(
		Yellow+"\n"+printFormat+Reset,
		playlist.Name,
		playlist.Contacts,
		playlist.FollowersTotal,
		strings.Join(firstThree(playlist.Styles), ", "),
		playlist.Region,
	)
}

func (s *SpotifyHandler) ProgressBar(pCh spotifind.ProgressChan) {
	startTime := time.Now()

	for p := range pCh {
		elapsed := time.Since(startTime)

		if p.Done == 0 {
			p.Done = 1
		}
		if p.Done > p.Total {
			p.Done = p.Total
		}
		averageTimePerUnit := elapsed / time.Duration(p.Done)
		remainingUnits := p.Total - p.Done
		estimatedTimeLeft := averageTimePerUnit * time.Duration(remainingUnits)

		percentage := int(math.Round((float64(p.Done) / float64(p.Total)) * 100))
		rootCmd.Printf(
			Red+"\r%d/%d - %d%% - time left: %s"+Reset,
			p.Done,
			p.Total,
			percentage,
			shortDur(estimatedTimeLeft),
		)
	}
}

func (s *SpotifyHandler) SearchPlaylistAllMarkets(q, ignore []string) {
	ch := make(spotifind.SpotifindChan)
	pCh := make(spotifind.ProgressChan)

	go func() {
		for i := 0; i < len(configs); i++ {
			err := s.spotifind.SearchPlaylistAllMarkets(ch, pCh, q, ignore)
			if errors.Is(err, spotifind.ErrTimeout) || errors.Is(err, spotifind.ErrTokenExpired) {
				if err = s.Reconnect(); err != nil {
					errHandler(err)
				}
			} else if err != nil {
				errHandler(err)
			}
		}
	}()
	go s.ProgressBar(pCh)
	for playlist := range ch {
		s.OutputPlaylist(playlist.Playlist)
	}
}

func (s *SpotifyHandler) SearchPlaylistForMarket(market string, q, ignore []string) {
	ch := make(spotifind.SpotifindChan)
	pCh := make(spotifind.ProgressChan)

	go func() {
		for i := 0; i < len(configs); i++ {
			err := s.spotifind.SearchPlaylistForMarket(ch, pCh, market, q, ignore)
			if errors.Is(err, spotifind.ErrTimeout) || errors.Is(err, spotifind.ErrTokenExpired) {
				if err = s.Reconnect(); err != nil {
					errHandler(err)
				}
			} else if err != nil {
				errHandler(err)
			}
		}
	}()
	go s.ProgressBar(pCh)
	for playlist := range ch {
		s.OutputPlaylist(playlist.Playlist)
	}
}

func (s *SpotifyHandler) SearchPlaylistPopular(q, ignore []string) {
	ch := make(spotifind.SpotifindChan)
	pCh := make(spotifind.ProgressChan)

	go func() {
		for i := 0; i < len(configs); i++ {
			err := s.spotifind.SearchPlaylistPopular(ch, pCh, q, ignore)
			if errors.Is(err, spotifind.ErrTimeout) || errors.Is(err, spotifind.ErrTokenExpired) {
				if err = s.Reconnect(); err != nil {
					errHandler(err)
				}
			} else if err != nil {
				errHandler(err)
			}
		}
	}()
	go s.ProgressBar(pCh)
	for playlist := range ch {
		s.OutputPlaylist(playlist.Playlist)
	}
}

func (s *SpotifyHandler) SearchPlaylistUnpopular(q, ignore []string) {
	ch := make(spotifind.SpotifindChan)
	pCh := make(spotifind.ProgressChan)

	go func() {
		for i := 0; i < len(configs); i++ {
			err := s.spotifind.SearchPlaylistUnpopular(ch, pCh, q, ignore)
			if errors.Is(err, spotifind.ErrTimeout) || errors.Is(err, spotifind.ErrTokenExpired) {
				if err = s.Reconnect(); err != nil {
					errHandler(err)
				}
			} else if err != nil {
				errHandler(err)
			}
		}
	}()
	go s.ProgressBar(pCh)
	for playlist := range ch {
		s.OutputPlaylist(playlist.Playlist)
	}
}

func (s *SpotifyHandler) OutputPlaylist(playlist spotifind.Playlist) {
	if s.IsPlaylistKnown(playlist.ExternalURLs["spotify"]) {
		return
	}
	s.PrintFormattedPlaylist(playlist)

	if err := s.Csv.WriteToFile(playlist); err != nil {
		fmt.Errorf("error while writing to file: %v", err)
	}
}

func (s *SpotifyHandler) IsPlaylistKnown(externalURL string) bool {
	if _, ok := s.KnownPlaylists[externalURL]; ok {
		return true
	}

	return false
}
