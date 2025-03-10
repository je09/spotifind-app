package csv

// that's a shit code for gui and cli applications, don't reuse it ever.

import (
	"encoding/csv"
	"fmt"
	"github.com/je09/spotifind"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type ReaderWriter interface {
	WriteToFile(playlist spotifind.Playlist) error
	ReadFromFile() (map[string]struct{} /* playlist uris */, error)
	SetFilePath(filePath string) error
}

type CSV struct {
	Path string
}

func New(path string) ReaderWriter {
	return &CSV{
		Path: path,
	}
}

func (c *CSV) WriteToFile(playlist spotifind.Playlist) error {
	if c.Path == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(c.Path), os.ModePerm); err != nil {
		return err
	}

	file, err := os.OpenFile(c.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 5 styles max
	if len(playlist.Styles) > 5 {
		playlist.Styles = playlist.Styles[:5]
	}

	// Convert the playlist to a slice of strings
	data := []string{
		playlist.Name,
		strconv.Itoa(playlist.FollowersTotal),
		strings.Join(playlist.Styles, "; "),
		strings.Join(playlist.Contacts, "; "),
		playlist.Description,
		playlist.Region,
		playlist.ExternalURLs["spotify"],
	}

	// remove all the comas
	for i, d := range data {
		data[i] = strings.ReplaceAll(d, ",", "")
	}

	if err := writer.Write(data); err != nil {
		return err
	}

	return nil
}

func (c *CSV) ReadFromFile() (map[string]struct{}, error) {
	file, err := os.Open(c.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	playlists := make(map[string]struct{})

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		playlist := spotifind.Playlist{
			Name: record[0],
			ExternalURLs: map[string]string{
				"spotify": record[6],
			},
		}

		playlists[playlist.ExternalURLs["spotify"]] = struct{}{}
	}

	return playlists, nil
}

func (c *CSV) SetFilePath(filePath string) error {
	if c.Path == "" {
		return fmt.Errorf("no save location set")
	}

	switch runtime.GOOS {
	case "windows":
		c.Path = fmt.Sprintf("%s\\%s.csv", c.Path, filePath)
	default:
		c.Path = fmt.Sprintf("%s/%s.csv", c.Path, filePath)
	}

	return nil
}
