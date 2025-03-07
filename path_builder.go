package main

import (
	"fmt"
	"os"
	"runtime"
)

var home, _ = os.UserHomeDir()

var (
	CacheLocationNix = fmt.Sprintf("%s/.spotifind/previous_searches.yml", home)
	CacheLocationWin = fmt.Sprintf("%s\\AppData\\Roaming\\spotifind\\previous_searches.yml", home)
)

var (
	ConfigLocationsDefault = []string{
		".",
	}

	ConfigLocationsNix = []string{
		fmt.Sprintf("%s/.config/spotifind", home),
		fmt.Sprintf("%s/.spotifind", home),
		fmt.Sprintf("%s/spotifind", home),
		fmt.Sprintf("%s/Spotifind", home),
		fmt.Sprintf("%s/Documents/spotifind", home),
		fmt.Sprintf("%s/Documents/Spotifind", home),
		"/etc/spotifind",
	}

	ConfigLocationsWin = []string{
		fmt.Sprintf("%s\\Documents\\spotifind", home),
		fmt.Sprintf("%s\\Documents\\Spotifind", home),
		fmt.Sprintf("%s\\AppData\\Roaming\\spotifind", home),
		"C:\\ProgramData\\spotifind",
	}
)

var (
	LogLocationDarwin = fmt.Sprintf("%s/Library/Logs/Spotifind/spotifind.log", home)
	LogLocationLinux  = fmt.Sprintf("%s/.spotifind/spotifind.log", home)
	LogLocationWin    = fmt.Sprintf("%s\\AppData\\Roaming\\spotifind\\spotifind.log", home)
)

type PathBuilder interface {
	ConfigLocations() []string
	CacheLocation() string
	LogLocations() string
}

type PathBuilderImpl struct {
	os string
}

func NewPathBuilder() PathBuilder {
	return &PathBuilderImpl{
		os: runtime.GOOS,
	}
}

func (p *PathBuilderImpl) ConfigLocations() []string {
	var locations []string
	locations = append(locations, ConfigLocationsDefault...)

	switch p.os {
	case "windows":
		locations = append(locations, ConfigLocationsWin...)
	default:
		locations = append(locations, ConfigLocationsNix...)
	}

	return locations
}

func (p *PathBuilderImpl) CacheLocation() string {
	switch p.os {
	case "windows":
		return CacheLocationWin
	default:
		return CacheLocationNix
	}
}

func (p *PathBuilderImpl) LogLocations() string {
	switch p.os {
	case "darwin":
		return LogLocationDarwin
	case "windows":
		return LogLocationWin
	default:
		return LogLocationLinux
	}
}
