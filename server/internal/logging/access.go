package logging

import (
	"net"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"hostbin/internal/router"
)

func Middleware(logger *zap.Logger, trustProxyHeaders bool, trustedProxyNets []*net.IPNet) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			capturingWriter := &statusWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(capturingWriter, r)

			meta := router.RequestMetaFromContext(r.Context())
			fields := []zap.Field{
				zap.String("request_id", metaValue(meta, func(m *router.RequestMeta) string { return m.RequestID })),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("host", hostForLog(meta, r.Host)),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", capturingWriter.status),
				zap.Int64("duration_ms", time.Since(started).Milliseconds()),
			}

			if meta != nil && meta.Slug != "" {
				fields = append(fields, zap.String("slug", meta.Slug))
			}

			if clientIP := forwardedClientIP(r, trustProxyHeaders, trustedProxyNets); clientIP != "" {
				fields = append(fields, zap.String("client_ip", clientIP))
			}

			logger.Info("request completed", fields...)
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func hostForLog(meta *router.RequestMeta, rawHost string) string {
	if meta != nil && meta.Host != "" {
		return meta.Host
	}
	return rawHost
}

func metaValue(meta *router.RequestMeta, extract func(*router.RequestMeta) string) string {
	if meta == nil {
		return ""
	}
	return extract(meta)
}

func forwardedClientIP(r *http.Request, trustProxyHeaders bool, trustedProxyNets []*net.IPNet) string {
	if !trustProxyHeaders || !remoteAddrTrusted(r.RemoteAddr, trustedProxyNets) {
		return ""
	}

	xff := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip := net.ParseIP(xff); ip != nil {
		return ip.String()
	}

	if ip := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP"))); ip != nil {
		return ip.String()
	}

	return ""
}

func remoteAddrTrusted(remoteAddr string, trustedProxyNets []*net.IPNet) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		host = remoteAddr
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	for _, network := range trustedProxyNets {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
