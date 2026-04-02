package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListenAddr         string
	BaseDomain         string
	AdminHost          string
	PresharedKey       string
	DBPath             string
	ReservedSubdomains []string
	ReservedSet        map[string]struct{}
	MaxDocSize         int64
	AuthTimestampSkew  time.Duration
	NonceTTL           time.Duration
	TrustProxyHeaders  bool
	TrustedProxyCIDRs  []string
	TrustedProxyNets   []*net.IPNet
	LogLevel           string
}

func Load() (Config, error) {
	maxDocSize, err := envInt64OrDefault("MAX_DOC_SIZE", 1_048_576)
	if err != nil {
		return Config{}, err
	}

	authTimestampSkewSeconds, err := envInt64OrDefault("AUTH_TIMESTAMP_SKEW_SECONDS", 60)
	if err != nil {
		return Config{}, err
	}

	nonceTTLSeconds, err := envInt64OrDefault("NONCE_TTL_SECONDS", 300)
	if err != nil {
		return Config{}, err
	}

	trustProxyHeaders, err := envBoolOrDefault("TRUST_PROXY_HEADERS", false)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		ListenAddr:         envOrDefault("LISTEN_ADDR", "127.0.0.1:8080"),
		BaseDomain:         strings.ToLower(strings.TrimSpace(os.Getenv("BASE_DOMAIN"))),
		AdminHost:          strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_HOST"))),
		PresharedKey:       strings.TrimSpace(os.Getenv("PRESHARED_KEY")),
		DBPath:             strings.TrimSpace(os.Getenv("DB_PATH")),
		ReservedSubdomains: splitCSV(envOrDefault("RESERVED_SUBDOMAINS", "admin,www,api")),
		MaxDocSize:         maxDocSize,
		AuthTimestampSkew:  time.Duration(authTimestampSkewSeconds) * time.Second,
		NonceTTL:           time.Duration(nonceTTLSeconds) * time.Second,
		TrustProxyHeaders:  trustProxyHeaders,
		TrustedProxyCIDRs:  splitCSV(envOrDefault("TRUSTED_PROXY_CIDRS", "127.0.0.1/32,::1/128")),
		LogLevel:           strings.ToLower(strings.TrimSpace(envOrDefault("LOG_LEVEL", "info"))),
	}

	reservedSet := make(map[string]struct{}, len(cfg.ReservedSubdomains))
	for _, entry := range cfg.ReservedSubdomains {
		reservedSet[entry] = struct{}{}
	}
	cfg.ReservedSet = reservedSet

	trustedProxyNets := make([]*net.IPNet, 0, len(cfg.TrustedProxyCIDRs))
	for _, cidr := range cfg.TrustedProxyCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return Config{}, fmt.Errorf("parse TRUSTED_PROXY_CIDRS %q: %w", cidr, err)
		}
		trustedProxyNets = append(trustedProxyNets, network)
	}
	cfg.TrustedProxyNets = trustedProxyNets

	if err := Validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envInt64OrDefault(key string, fallback int64) (int64, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}

	return parsed, nil
}

func envBoolOrDefault(key string, fallback bool) (bool, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("parse %s: %w", key, err)
	}

	return parsed, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.ToLower(strings.TrimSpace(part))
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
