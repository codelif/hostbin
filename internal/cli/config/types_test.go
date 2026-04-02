package config

import "testing"

func TestValidateRejectsServerURLWithPath(t *testing.T) {
	cfg := File{
		ServerURL: "https://admin.domain.com/api/v1",
		AuthKey:   "01234567890123456789012345678901",
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want error")
	}
}

func TestValidatePartialAllowsIncrementalConfig(t *testing.T) {
	tests := []File{
		{ServerURL: "https://admin.domain.com"},
		{AuthKey: "01234567890123456789012345678901"},
		{Timeout: "5s"},
	}

	for _, cfg := range tests {
		if err := cfg.ValidatePartial(); err != nil {
			t.Fatalf("ValidatePartial(%+v) error = %v", cfg, err)
		}
	}
}
