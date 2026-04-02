package config

import (
	"strings"
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{
			name: "invalid listen addr",
			mutate: func(cfg *Config) {
				cfg.ListenAddr = "bad-addr"
			},
			wantErr: "invalid LISTEN_ADDR",
		},
		{
			name: "invalid base domain",
			mutate: func(cfg *Config) {
				cfg.BaseDomain = "bad_domain"
			},
			wantErr: "invalid BASE_DOMAIN",
		},
		{
			name: "admin host not under base domain",
			mutate: func(cfg *Config) {
				cfg.AdminHost = "admin.other.com"
			},
			wantErr: "ADMIN_HOST must be under BASE_DOMAIN",
		},
		{
			name: "short preshared key",
			mutate: func(cfg *Config) {
				cfg.PresharedKey = "short"
			},
			wantErr: "PRESHARED_KEY shared secret must be at least 32 bytes",
		},
		{
			name: "empty db path",
			mutate: func(cfg *Config) {
				cfg.DBPath = ""
			},
			wantErr: "DB_PATH must not be empty",
		},
		{
			name: "non-positive max doc size",
			mutate: func(cfg *Config) {
				cfg.MaxDocSize = 0
			},
			wantErr: "MAX_DOC_SIZE must be greater than zero",
		},
		{
			name: "non-positive timestamp skew",
			mutate: func(cfg *Config) {
				cfg.AuthTimestampSkew = 0
			},
			wantErr: "AUTH_TIMESTAMP_SKEW_SECONDS must be greater than zero",
		},
		{
			name: "non-positive nonce ttl",
			mutate: func(cfg *Config) {
				cfg.NonceTTL = 0
			},
			wantErr: "NONCE_TTL_SECONDS must be greater than zero",
		},
		{
			name: "invalid reserved subdomain",
			mutate: func(cfg *Config) {
				cfg.ReservedSubdomains = []string{"bad_slug"}
			},
			wantErr: "invalid reserved subdomain",
		},
		{
			name: "invalid log level",
			mutate: func(cfg *Config) {
				cfg.LogLevel = "verbose"
			},
			wantErr: "invalid LOG_LEVEL",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validConfig()
			tc.mutate(&cfg)

			err := Validate(cfg)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Validate() error = %v, want substring %q", err, tc.wantErr)
			}
		})
	}
}

func validConfig() Config {
	return Config{
		ListenAddr:         "127.0.0.1:8080",
		BaseDomain:         "domain.com",
		AdminHost:          "admin.domain.com",
		PresharedKey:       "01234567890123456789012345678901",
		DBPath:             "./data.db",
		ReservedSubdomains: []string{"admin", "www", "api"},
		MaxDocSize:         1_048_576,
		AuthTimestampSkew:  60 * time.Second,
		NonceTTL:           300 * time.Second,
		LogLevel:           "info",
	}
}
