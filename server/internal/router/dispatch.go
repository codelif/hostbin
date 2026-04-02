package router

import (
	"net/http"
)

type Dispatcher struct {
	baseDomain string
	adminHost  string
	reserved   map[string]struct{}
	admin      http.Handler
	public     http.Handler
}

func NewDispatcher(baseDomain, adminHost string, reserved map[string]struct{}, admin, public http.Handler) *Dispatcher {
	return &Dispatcher{
		baseDomain: baseDomain,
		adminHost:  adminHost,
		reserved:   reserved,
		admin:      admin,
		public:     public,
	}
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info, err := ClassifyHost(r.Host, d.baseDomain, d.adminHost, d.reserved)
	if meta := RequestMetaFromContext(r.Context()); meta != nil {
		meta.Host = info.Host
		meta.HostKind = string(info.Kind)
		meta.Slug = info.Slug
	}

	if err != nil {
		writePlaintext(w, http.StatusNotFound, "not found\n")
		return
	}

	switch info.Kind {
	case HostKindAdmin:
		d.admin.ServeHTTP(w, r)
	case HostKindPublic:
		d.public.ServeHTTP(w, r)
	default:
		writePlaintext(w, http.StatusNotFound, "not found\n")
	}
}

func writePlaintext(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
