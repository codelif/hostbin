package httpserver

import (
	"net/http"

	"go.uber.org/zap"

	"hostbin/internal/router"
)

func Recovery(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				meta := router.RequestMetaFromContext(r.Context())
				logger.Error("panic recovered",
					zap.Any("panic", recovered),
					zap.String("request_id", metaValue(meta, func(m *router.RequestMeta) string { return m.RequestID })),
					zap.String("host", metaValue(meta, func(m *router.RequestMeta) string { return m.Host })),
					zap.String("path", r.URL.Path),
				)

				if meta != nil && meta.HostKind == string(router.HostKindAdmin) {
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

func metaValue(meta *router.RequestMeta, extract func(*router.RequestMeta) string) string {
	if meta == nil {
		return ""
	}
	return extract(meta)
}
