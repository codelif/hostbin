package config

import (
	"fmt"
	"net"
	"strings"

	"hostbin/internal/domain/slugs"
)

func Validate(cfg Config) error {
	if _, err := net.ResolveTCPAddr("tcp", cfg.ListenAddr); err != nil {
		return fmt.Errorf("invalid LISTEN_ADDR: %w", err)
	}

	if !validHostname(cfg.BaseDomain) {
		return fmt.Errorf("invalid BASE_DOMAIN")
	}

	if !validHostname(cfg.AdminHost) {
		return fmt.Errorf("invalid ADMIN_HOST")
	}

	if cfg.AdminHost == cfg.BaseDomain || !strings.HasSuffix(cfg.AdminHost, "."+cfg.BaseDomain) {
		return fmt.Errorf("ADMIN_HOST must be under BASE_DOMAIN")
	}

	if len(cfg.PresharedKey) < 32 {
		return fmt.Errorf("PRESHARED_KEY must be at least 32 bytes")
	}

	if strings.TrimSpace(cfg.DBPath) == "" {
		return fmt.Errorf("DB_PATH must not be empty")
	}

	if cfg.MaxDocSize <= 0 {
		return fmt.Errorf("MAX_DOC_SIZE must be greater than zero")
	}

	if cfg.AuthTimestampSkew <= 0 {
		return fmt.Errorf("AUTH_TIMESTAMP_SKEW_SECONDS must be greater than zero")
	}

	if cfg.NonceTTL <= 0 {
		return fmt.Errorf("NONCE_TTL_SECONDS must be greater than zero")
	}

	for _, entry := range cfg.ReservedSubdomains {
		if err := slugs.Validate(entry, nil); err != nil {
			return fmt.Errorf("invalid reserved subdomain %q", entry)
		}
	}

	switch cfg.LogLevel {
	case "debug", "info", "warn", "error", "dpanic", "panic", "fatal":
	default:
		return fmt.Errorf("invalid LOG_LEVEL")
	}

	return nil
}

func validHostname(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) < 2 {
		return false
	}

	for _, part := range parts {
		if err := slugs.Validate(part, nil); err != nil {
			return false
		}
	}

	return true
}
