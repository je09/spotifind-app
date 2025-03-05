package main

import (
	"github.com/je09/spotifind"
	"github.com/je09/spotifind-app/pkg/csv"
)

type CsvHandlerImpl struct {
	handler *csv.CsvHandler
}

func NewCsvHandler(path string) *CsvHandlerImpl {
	return &CsvHandlerImpl{
		handler: csv.NewCsvHandler(path),
	}
}

func (c *CsvHandlerImpl) ReadFromFile() ([]string, error) {
	return c.handler.ReadFromFile()
}

func (c *CsvHandlerImpl) WriteToFile(playlist spotifind.Playlist) error {
	return c.handler.WriteToFile(playlist)
}
