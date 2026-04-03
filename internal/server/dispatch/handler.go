package dispatch

import (
	"net/http"

	"github.com/codelif/hostbin/internal/domain/hosts"
	"github.com/codelif/hostbin/internal/server/requestmeta"
)

type Dispatcher struct {
	baseDomain string
	adminHost  string
	reserved   map[string]struct{}
	admin      http.Handler
	public     http.Handler
}

func NewHandler(baseDomain, adminHost string, reserved map[string]struct{}, admin, public http.Handler) *Dispatcher {
	return &Dispatcher{
		baseDomain: baseDomain,
		adminHost:  adminHost,
		reserved:   reserved,
		admin:      admin,
		public:     public,
	}
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info, err := hosts.ClassifyHost(r.Host, d.baseDomain, d.adminHost, d.reserved)
	if meta := requestmeta.FromContext(r.Context()); meta != nil {
		meta.Host = info.Host
		meta.HostKind = info.Kind
		meta.Slug = info.Slug
	}

	if err != nil {
		writePlaintext(w, http.StatusNotFound, "not found\n")
		return
	}

	switch info.Kind {
	case hosts.KindAdmin:
		d.admin.ServeHTTP(w, r)
	case hosts.KindPublic:
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
