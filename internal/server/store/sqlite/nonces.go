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

package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/codelif/hostbin/internal/server/nonce"
)

type NonceStore struct {
	db  *sql.DB
	ttl time.Duration
}

func NewNonceStore(db *sql.DB, ttl time.Duration) *NonceStore {
	return &NonceStore{db: db, ttl: ttl}
}

func (s *NonceStore) UseOnce(value string, now time.Time) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := cleanupExpiredNonces(tx, now.UTC()); err != nil {
		return err
	}

	const query = `INSERT INTO nonces (nonce, expires_at) VALUES (?, ?)`
	_, err = tx.ExecContext(context.Background(), query, value, now.UTC().Add(s.ttl).Format(time.RFC3339))
	if err != nil {
		if isUniqueConstraint(err) {
			return nonce.ErrReplayed
		}
		return err
	}

	return tx.Commit()
}

func (s *NonceStore) CleanupExpired(now time.Time) {
	_ = cleanupExpiredNonces(s.db, now.UTC())
}

func cleanupExpiredNonces(execer interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, now time.Time) error {
	const query = `DELETE FROM nonces WHERE expires_at <= ?`
	_, err := execer.ExecContext(context.Background(), query, now.Format(time.RFC3339))
	return err
}
