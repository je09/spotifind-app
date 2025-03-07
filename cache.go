package main

import (
	"fmt"
	"github.com/je09/spotifind-app/common"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type PreviousSearch struct {
	Searches []string `yaml:"searches"`
	Ignores  []string `yaml:"ignores"`
}

type Cache interface {
	Load() error
	PreviousSearch() PreviousSearch
	Append(search, ignore string) error
}

type CacheImpl struct {
	path string
	prev PreviousSearch
}

func NewCache() *CacheImpl {
	pb := common.NewPathBuilder()
	return &CacheImpl{
		path: pb.CacheLocation(),
	}
}

func (c *CacheImpl) Load() error {
	// if file does not exist, create it
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		_ = c.Save()
	}

	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}

	var prev PreviousSearch
	if err := yaml.Unmarshal(data, &prev); err != nil {
		return err
	}

	c.prev = prev

	return nil
}

func (c *CacheImpl) PreviousSearch() PreviousSearch {
	return c.prev
}

func (c *CacheImpl) Append(search, ignore string) error {
	if err := c.AppendSearch(search); err != nil {
		return err
	}
	if err := c.AppendIgnore(ignore); err != nil {
		return err
	}

	if err := c.Save(); err != nil {
		return err
	}

	return nil
}

func (c *CacheImpl) AppendSearch(search string) error {
	for _, s := range c.prev.Searches {
		if s == search {
			return nil
		}
	}

	c.prev.Searches = append(c.prev.Searches, search)
	if len(c.prev.Searches) > 5 {
		c.prev.Searches = c.prev.Searches[1:]
	}

	return nil
}

func (c *CacheImpl) AppendIgnore(ignore string) error {
	for _, i := range c.prev.Ignores {
		if i == ignore {
			return nil
		}
	}

	c.prev.Ignores = append(c.prev.Ignores, ignore)
	if len(c.prev.Ignores) > 5 {
		c.prev.Ignores = c.prev.Ignores[1:]
	}

	return nil
}

func (c *CacheImpl) Save() error {
	// check if root
	if os.Geteuid() == 0 {
		fmt.Printf("Running as root is not recommended.")
	}

	if err := os.MkdirAll(filepath.Dir(c.path), os.ModePerm); err != nil {
		return err
	}

	prevMarshalled, err := yaml.Marshal(c.prev)
	if err != nil {
		return err
	}

	if err := os.WriteFile(c.path, prevMarshalled, 0644); err != nil {
		return err
	}

	return nil
}
