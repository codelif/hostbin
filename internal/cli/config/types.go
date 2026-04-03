// Copyright (c) 2026 Harsh Sharma <harsh@codelif.in>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/codelif/hostbin/internal/protocol/authsig"
)

const (
	DefaultServerURL = "https://admin.domain.com"
	DefaultTimeout   = "10s"
	DefaultColor     = "auto"
)

type File struct {
	ServerURL string `toml:"server_url"`
	AuthKey   string `toml:"auth_key"`
	Timeout   string `toml:"timeout"`
	Editor    string `toml:"editor"`
	Color     string `toml:"color"`
}

func DefaultFile() File {
	return File{
		Timeout: DefaultTimeout,
		Color:   DefaultColor,
	}
}

func (f File) Normalized() File {
	trimmed := File{
		ServerURL: strings.TrimSpace(f.ServerURL),
		AuthKey:   strings.TrimSpace(f.AuthKey),
		Timeout:   strings.TrimSpace(f.Timeout),
		Editor:    strings.TrimSpace(f.Editor),
		Color:     strings.TrimSpace(strings.ToLower(f.Color)),
	}

	trimmed.ServerURL = strings.TrimRight(trimmed.ServerURL, "/")

	if trimmed.Timeout == "" {
		trimmed.Timeout = DefaultTimeout
	}
	if trimmed.Color == "" {
		trimmed.Color = DefaultColor
	}

	return trimmed
}

func (f File) Validate() error {
	normalized := f.Normalized()
	if err := normalized.ValidatePartial(); err != nil {
		return err
	}

	if normalized.ServerURL == "" {
		return fmt.Errorf("server_url is required")
	}
	if normalized.AuthKey == "" {
		return fmt.Errorf("auth_key is required")
	}

	return nil
}

func (f File) ValidatePartial() error {
	normalized := f.Normalized()
	if normalized.ServerURL != "" {
		parsed, err := url.Parse(normalized.ServerURL)
		if err != nil {
			return fmt.Errorf("invalid server_url: %w", err)
		}
		if !parsed.IsAbs() || parsed.Host == "" {
			return fmt.Errorf("server_url must be an absolute URL")
		}
		if parsed.Scheme != "https" && parsed.Scheme != "http" {
			return fmt.Errorf("server_url must use http or https")
		}
		if parsed.Path != "" && parsed.Path != "/" {
			return fmt.Errorf("server_url must not include a path")
		}
		if parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("server_url must not include query or fragment")
		}
	}

	if normalized.AuthKey != "" {
		if err := authsig.ValidateSharedSecret(normalized.AuthKey); err != nil {
			return fmt.Errorf("auth_key %w", err)
		}
	}

	if _, err := time.ParseDuration(normalized.Timeout); err != nil {
		return fmt.Errorf("invalid timeout: %w", err)
	}

	switch normalized.Color {
	case "auto", "always", "never":
	default:
		return fmt.Errorf("color must be auto, always, or never")
	}

	return nil
}

func (f File) Duration() (time.Duration, error) {
	return time.ParseDuration(f.Normalized().Timeout)
}

func (f File) Get(key string) (string, error) {
	normalized := f.Normalized()
	switch key {
	case "server_url":
		return normalized.ServerURL, nil
	case "auth_key":
		return normalized.AuthKey, nil
	case "timeout":
		return normalized.Timeout, nil
	case "editor":
		return normalized.Editor, nil
	case "color":
		return normalized.Color, nil
	default:
		return "", fmt.Errorf("unknown config key %q", key)
	}
}

func (f *File) Set(key, value string) error {
	switch key {
	case "server_url":
		f.ServerURL = value
	case "auth_key":
		f.AuthKey = value
	case "timeout":
		f.Timeout = value
	case "editor":
		f.Editor = value
	case "color":
		f.Color = value
	default:
		return fmt.Errorf("unknown config key %q", key)
	}

	*f = f.Normalized()
	return f.ValidatePartial()
}

func Keys() []string {
	return []string{"server_url", "auth_key", "timeout", "editor", "color"}
}

func MaskSecret(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return "************"
}
