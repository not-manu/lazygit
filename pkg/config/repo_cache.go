package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type repoCache[T any] struct {
	mutex       sync.Mutex
	path        string
	description string
	byRepoPath  map[string]T
	loadErr     error
}

func loadRepoCache[T any](fileName string, description string) *repoCache[T] {
	path, err := repoCachePath(fileName)
	if err != nil {
		cache := newRepoCache[T]("", description)
		cache.loadErr = err
		return cache
	}

	cache := newRepoCache[T](path, description)
	cache.load()
	return cache
}

func repoCachePath(fileName string) (string, error) {
	path, err := stateFilePath(stateFileName)
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(path), fileName), nil
}

func newRepoCache[T any](path string, description string) *repoCache[T] {
	return &repoCache[T]{
		path:        path,
		description: description,
		byRepoPath:  make(map[string]T),
	}
}

func (c *repoCache[T]) load() {
	if c.path == "" {
		return
	}

	content, err := os.ReadFile(c.path)
	if err != nil {
		if !os.IsNotExist(err) {
			c.loadErr = fmt.Errorf("reading %s cache: %w", c.description, err)
		}
		return
	}
	if len(content) == 0 {
		return
	}

	if err := json.Unmarshal(content, &c.byRepoPath); err != nil {
		c.byRepoPath = make(map[string]T)
		c.loadErr = fmt.Errorf("parsing %s cache: %w", c.description, err)
	} else if c.byRepoPath == nil {
		c.byRepoPath = make(map[string]T)
	}
}

func (c *repoCache[T]) get(repoPath string) T {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.byRepoPath[repoPath]
}

func (c *repoCache[T]) takeLoadError() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	loadErr := c.loadErr
	c.loadErr = nil
	return loadErr
}

func (c *repoCache[T]) save(repoPath string, value T) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.byRepoPath[repoPath] = value
	if c.path == "" {
		return nil
	}

	content, err := json.MarshalIndent(c.byRepoPath, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')

	if err := os.WriteFile(c.path, content, 0o644); err != nil && !os.IsPermission(err) {
		return err
	}

	return nil
}
