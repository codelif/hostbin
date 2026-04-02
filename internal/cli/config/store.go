package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

var ErrNotFound = errors.New("config not found")

type Store struct {
	path string
}

func NewStore(explicitPath string) (*Store, error) {
	path := explicitPath
	if path == "" {
		var err error
		path, err = DefaultPath()
		if err != nil {
			return nil, err
		}
	}

	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	return &Store{path: absolute}, nil
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) Load() (File, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return File{}, ErrNotFound
		}
		return File{}, err
	}

	var cfg File
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return File{}, fmt.Errorf("parse config: %w", err)
	}

	cfg = cfg.Normalized()
	if err := cfg.ValidatePartial(); err != nil {
		return File{}, err
	}

	return cfg, nil
}

func (s *Store) Save(cfg File) error {
	cfg = cfg.Normalized()
	if err := cfg.ValidatePartial(); err != nil {
		return err
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	if err := os.Chmod(dir, 0o700); err != nil && !errors.Is(err, os.ErrPermission) {
		return err
	}

	tempFile, err := os.CreateTemp(dir, "hbcli-config-*.toml")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer func() { _ = os.Remove(tempPath) }()

	if err := tempFile.Chmod(0o600); err != nil {
		_ = tempFile.Close()
		return err
	}
	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempPath, s.path); err != nil {
		return err
	}

	return os.Chmod(s.path, 0o600)
}
