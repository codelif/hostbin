package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadDefaultsAndNormalization(t *testing.T) {
	setValidEnv(t)
	t.Setenv("BASE_DOMAIN", "Domain.COM")
	t.Setenv("ADMIN_HOST", "Admin.Domain.COM")
	t.Setenv("RESERVED_SUBDOMAINS", "admin, WWW, api")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.ListenAddr != "127.0.0.1:8080" {
		t.Fatalf("ListenAddr = %q, want %q", cfg.ListenAddr, "127.0.0.1:8080")
	}
	if cfg.BaseDomain != "domain.com" {
		t.Fatalf("BaseDomain = %q, want %q", cfg.BaseDomain, "domain.com")
	}
	if cfg.AdminHost != "admin.domain.com" {
		t.Fatalf("AdminHost = %q, want %q", cfg.AdminHost, "admin.domain.com")
	}
	if cfg.MaxDocSize != 1_048_576 {
		t.Fatalf("MaxDocSize = %d, want %d", cfg.MaxDocSize, 1_048_576)
	}
	if cfg.AuthTimestampSkew != 60*time.Second {
		t.Fatalf("AuthTimestampSkew = %s, want %s", cfg.AuthTimestampSkew, 60*time.Second)
	}
	if cfg.NonceTTL != 300*time.Second {
		t.Fatalf("NonceTTL = %s, want %s", cfg.NonceTTL, 300*time.Second)
	}
	if cfg.TrustProxyHeaders {
		t.Fatal("TrustProxyHeaders = true, want false")
	}
	if len(cfg.TrustedProxyNets) != 2 {
		t.Fatalf("len(TrustedProxyNets) = %d, want 2", len(cfg.TrustedProxyNets))
	}
	if _, ok := cfg.ReservedSet["www"]; !ok {
		t.Fatal("expected reserved set to include www")
	}
}

func TestLoadRejectsInvalidNumericEnv(t *testing.T) {
	setValidEnv(t)
	t.Setenv("MAX_DOC_SIZE", "nope")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "parse MAX_DOC_SIZE") {
		t.Fatalf("Load() error = %v, want parse MAX_DOC_SIZE", err)
	}
}

func TestLoadRejectsInvalidBooleanEnv(t *testing.T) {
	setValidEnv(t)
	t.Setenv("TRUST_PROXY_HEADERS", "sometimes")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "parse TRUST_PROXY_HEADERS") {
		t.Fatalf("Load() error = %v, want parse TRUST_PROXY_HEADERS", err)
	}
}

func TestLoadRejectsInvalidCIDR(t *testing.T) {
	setValidEnv(t)
	t.Setenv("TRUSTED_PROXY_CIDRS", "not-a-cidr")

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "parse TRUSTED_PROXY_CIDRS") {
		t.Fatalf("Load() error = %v, want parse TRUSTED_PROXY_CIDRS", err)
	}
}

func setValidEnv(t *testing.T) {
	t.Helper()
	t.Setenv("LISTEN_ADDR", "127.0.0.1:8080")
	t.Setenv("BASE_DOMAIN", "domain.com")
	t.Setenv("ADMIN_HOST", "admin.domain.com")
	t.Setenv("PRESHARED_KEY", "01234567890123456789012345678901")
	t.Setenv("DB_PATH", "./test.db")
	t.Setenv("RESERVED_SUBDOMAINS", "admin,www,api")
	t.Setenv("MAX_DOC_SIZE", "")
	t.Setenv("AUTH_TIMESTAMP_SKEW_SECONDS", "")
	t.Setenv("NONCE_TTL_SECONDS", "")
	t.Setenv("TRUST_PROXY_HEADERS", "")
	t.Setenv("TRUSTED_PROXY_CIDRS", "")
	t.Setenv("LOG_LEVEL", "")
}
