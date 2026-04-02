package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvLoadsMissingValuesOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("BASE_DOMAIN=example.com\nPRESHARED_KEY=from-file\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	t.Setenv("PRESHARED_KEY", "already-set")

	if err := LoadDotEnv(path); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}

	if got := os.Getenv("BASE_DOMAIN"); got != "example.com" {
		t.Fatalf("BASE_DOMAIN = %q, want %q", got, "example.com")
	}
	if got := os.Getenv("PRESHARED_KEY"); got != "already-set" {
		t.Fatalf("PRESHARED_KEY = %q, want %q", got, "already-set")
	}
}

func TestLoadDotEnvIgnoresMissingFile(t *testing.T) {
	if err := LoadDotEnv(filepath.Join(t.TempDir(), "missing.env")); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}
}
