package main

import (
	"context"
	"fmt"
	"github.com/je09/spotifind"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Search(q, ignore, market string) string {
	return fmt.Sprintf("Query: %s, market: %s, ignore: %s", q, market, ignore)
}

func (a *App) Markets() []string {
	return spotifind.MarketsAll
}
