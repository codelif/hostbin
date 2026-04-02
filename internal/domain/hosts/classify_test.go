package hosts

import "testing"

func TestClassifyHost(t *testing.T) {
	reserved := map[string]struct{}{"admin": {}, "www": {}}

	tests := []struct {
		name     string
		host     string
		wantKind Kind
		wantSlug string
		wantErr  bool
	}{
		{name: "public host", host: "doc1.domain.com", wantKind: KindPublic, wantSlug: "doc1"},
		{name: "public host with port", host: "doc1.domain.com:443", wantKind: KindPublic, wantSlug: "doc1"},
		{name: "admin host", host: "admin.domain.com", wantKind: KindAdmin},
		{name: "base domain only", host: "domain.com", wantErr: true},
		{name: "multi level subdomain", host: "a.b.domain.com", wantErr: true},
		{name: "reserved public slug", host: "www.domain.com", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			info, err := ClassifyHost(tc.host, "domain.com", "admin.domain.com", reserved)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for host %q", tc.host)
				}
				return
			}

			if err != nil {
				t.Fatalf("ClassifyHost(%q) error = %v", tc.host, err)
			}
			if info.Kind != tc.wantKind {
				t.Fatalf("kind = %q, want %q", info.Kind, tc.wantKind)
			}
			if info.Slug != tc.wantSlug {
				t.Fatalf("slug = %q, want %q", info.Slug, tc.wantSlug)
			}
		})
	}
}
