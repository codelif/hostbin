package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/codelif/hostbin/internal/domain/hosts"
	"github.com/codelif/hostbin/internal/server/requestmeta"
)

func Recovery(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				meta := requestmeta.FromContext(r.Context())
				logger.Error("panic recovered",
					zap.Any("panic", recovered),
					zap.String("request_id", metaValue(meta, func(m *requestmeta.Meta) string { return m.RequestID })),
					zap.String("host", metaValue(meta, func(m *requestmeta.Meta) string { return m.Host })),
					zap.String("path", r.URL.Path),
				)

				if meta != nil && meta.HostKind == hosts.KindAdmin {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(`{"error":"internal_error"}`))
					return
				}

				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("internal error\n"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func metaValue(meta *requestmeta.Meta, extract func(*requestmeta.Meta) string) string {
	if meta == nil {
		return ""
	}
	return extract(meta)
}
