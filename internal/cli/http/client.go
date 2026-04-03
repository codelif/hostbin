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

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	cliconfig "github.com/codelif/hostbin/internal/cli/config"
	"github.com/codelif/hostbin/internal/protocol/adminv1"
	"github.com/codelif/hostbin/internal/protocol/authsig"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	secret     []byte
}

type APIError struct {
	StatusCode int
	Code       string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return e.Code
	}
	return fmt.Sprintf("unexpected status %d", e.StatusCode)
}

func New(cfg cliconfig.File) (*Client, error) {
	normalized := cfg.Normalized()
	if err := normalized.Validate(); err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(normalized.ServerURL)
	if err != nil {
		return nil, err
	}

	timeout, err := normalized.Duration()
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		secret: []byte(normalized.AuthKey),
	}, nil
}

func (c *Client) Health(ctx context.Context) (*adminv1.StatusResponse, error) {
	var response adminv1.StatusResponse
	if err := c.doJSON(ctx, http.MethodGet, adminv1.HealthPath, nil, false, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) AuthCheck(ctx context.Context) (*adminv1.StatusResponse, error) {
	var response adminv1.StatusResponse
	if err := c.doJSON(ctx, http.MethodGet, adminv1.AuthCheckPath, nil, true, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) ListDocuments(ctx context.Context) ([]adminv1.DocumentResponse, error) {
	var response adminv1.ListDocumentsResponse
	if err := c.doJSON(ctx, http.MethodGet, adminv1.DocumentsCollection, nil, true, &response); err != nil {
		return nil, err
	}

	return response.Documents, nil
}

func (c *Client) GetDocument(ctx context.Context, slug string) (*adminv1.DocumentResponse, error) {
	var response adminv1.DocumentResponse
	if err := c.doJSON(ctx, http.MethodGet, adminv1.DocumentPath(slug), nil, true, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetDocumentContent(ctx context.Context, slug string) ([]byte, error) {
	return c.doBytes(ctx, http.MethodGet, adminv1.DocumentContentPath(slug), nil, true, "")
}

func (c *Client) CreateDocument(ctx context.Context, slug string, content []byte) (*adminv1.DocumentResponse, error) {
	var response adminv1.DocumentResponse
	if err := c.doJSON(ctx, http.MethodPost, adminv1.DocumentPath(slug), content, true, &response, withContentType("text/plain; charset=utf-8")); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) ReplaceDocument(ctx context.Context, slug string, content []byte) (*adminv1.DocumentResponse, error) {
	var response adminv1.DocumentResponse
	if err := c.doJSON(ctx, http.MethodPut, adminv1.DocumentPath(slug), content, true, &response, withContentType("text/plain; charset=utf-8")); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DeleteDocument(ctx context.Context, slug string) (*adminv1.DeleteResponse, error) {
	var response adminv1.DeleteResponse
	if err := c.doJSON(ctx, http.MethodDelete, adminv1.DocumentPath(slug), nil, true, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type requestOption func(*http.Request)

func withContentType(contentType string) requestOption {
	return func(r *http.Request) {
		r.Header.Set("Content-Type", contentType)
	}
}

func (c *Client) doJSON(ctx context.Context, method, path string, body []byte, signed bool, target any, options ...requestOption) error {
	responseBody, err := c.doBytes(ctx, method, path, body, signed, "application/json", options...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	if err := json.Unmarshal(responseBody, target); err != nil {
		return err
	}

	return nil
}

func (c *Client) doBytes(ctx context.Context, method, path string, body []byte, signed bool, accept string, options ...requestOption) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.resolve(path), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, option := range options {
		option(req)
	}

	if signed {
		if err := authsig.SetSignedHeaders(req, body, c.secret, time.Now().UTC()); err != nil {
			return nil, err
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResponse adminv1.ErrorResponse
		if err := json.Unmarshal(responseBody, &errorResponse); err == nil && errorResponse.Error != "" {
			return nil, &APIError{StatusCode: resp.StatusCode, Code: errorResponse.Error}
		}
		return nil, &APIError{StatusCode: resp.StatusCode}
	}

	return responseBody, nil
}

func (c *Client) resolve(path string) string {
	resolved := *c.baseURL
	resolved.Path = strings.TrimRight(c.baseURL.Path, "/") + path
	resolved.RawPath = resolved.Path
	return resolved.String()
}
