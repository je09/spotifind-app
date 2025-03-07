package main

import (
	"github.com/je09/spotifind-app/cli"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
