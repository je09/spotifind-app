package common

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
)

const (
	dev = "dev"
)

type ReleaseManager interface {
	NewRelease(v string) (string, error)
}

type ReleaseManagerImpl struct {
}

func NewReleaseManager() *ReleaseManagerImpl {
	return &ReleaseManagerImpl{}
}

func (r *ReleaseManagerImpl) NewRelease(v string) (string, error) {
	resp, err := http.Get("https://api.github.com/repos/je09/spotifind-app/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	log.Infof("Current version: %s", v)
	if release.TagName == v || v == dev {
		return "", nil
	}

	return release.TagName, nil
}
