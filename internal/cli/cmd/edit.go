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

package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	clieditor "github.com/codelif/hostbin/internal/cli/editor"
	"github.com/codelif/hostbin/internal/cli/ui"
	"github.com/codelif/hostbin/internal/protocol/adminv1"
)

func newEditCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "edit <slug>",
		Short: "Edit an existing document in your editor",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			if err := validateSlug(slug); err != nil {
				return err
			}

			_, cfg, client, err := loadClient(app)
			if err != nil {
				return err
			}

			var initial []byte
			err = ui.RunSpinner(app.Stderr, "Fetching "+slug, func() error {
				content, err := client.GetDocumentContent(cmd.Context(), slug)
				if err == nil {
					initial = content
				}
				return err
			})
			if err != nil {
				return fmt.Errorf("fetch document: %w", err)
			}

			updated, changed, err := clieditor.EditBuffer(clieditor.Resolve(cfg.Editor), slug, initial)
			if err != nil {
				return err
			}
			if !changed || bytes.Equal(initial, updated) {
				_, err = fmt.Fprintln(app.Stdout, "No changes")
				return err
			}

			var doc *adminv1.DocumentResponse
			err = ui.RunSpinner(app.Stderr, "Uploading "+slug, func() error {
				response, err := client.ReplaceDocument(cmd.Context(), slug, updated)
				if err == nil {
					doc = response
				}
				return err
			})
			if err != nil {
				return fmt.Errorf("update document: %w", err)
			}

			return ui.PrintDocumentSummary(app.Stdout, "Updated", *doc)
		},
	}
}
