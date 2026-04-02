package slugs

import "testing"

func TestValidate(t *testing.T) {
	reserved := map[string]struct{}{"admin": {}}

	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "simple", value: "doc1"},
		{name: "hyphen", value: "my-notes"},
		{name: "reserved", value: "admin", wantErr: ErrReserved},
		{name: "underscore", value: "bad_slug", wantErr: ErrInvalid},
		{name: "leading hyphen", value: "-bad", wantErr: ErrInvalid},
		{name: "trailing hyphen", value: "bad-", wantErr: ErrInvalid},
		{name: "uppercase", value: "Doc1", wantErr: ErrInvalid},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.value, reserved)
			if err != tc.wantErr {
				t.Fatalf("Validate(%q) error = %v, want %v", tc.value, err, tc.wantErr)
			}
		})
	}
}
