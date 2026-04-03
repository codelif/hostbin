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
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	clieditor "github.com/codelif/hostbin/internal/cli/editor"
	"github.com/codelif/hostbin/internal/cli/input"
	"github.com/codelif/hostbin/internal/cli/ui"
	"github.com/codelif/hostbin/internal/protocol/adminv1"
)

func newNewCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "new <slug> [file]",
		Short: "Create a new document",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWriteCommand(cmd, app, args, true)
		},
	}
}

func newPutCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "put <slug> [file]",
		Short: "Replace an existing document",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWriteCommand(cmd, app, args, false)
		},
	}
}

func runWriteCommand(cmd *cobra.Command, app *cliapp.App, args []string, create bool) error {
	slug := args[0]
	if err := validateSlug(slug); err != nil {
		return err
	}

	_, cfg, client, err := loadClient(app)
	if err != nil {
		return err
	}

	content, err := resolveWriteContent(cmd, cfg.Editor, slug, args)
	if err != nil {
		return err
	}

	label := "Uploading " + slug
	var docErr error
	var action string
	var responseDoc any
	err = ui.RunSpinner(app.Stderr, label, func() error {
		if create {
			doc, err := client.CreateDocument(cmd.Context(), slug, content)
			responseDoc = doc
			docErr = err
			action = "Created"
			return err
		}
		doc, err := client.ReplaceDocument(cmd.Context(), slug, content)
		responseDoc = doc
		docErr = err
		action = "Updated"
		return err
	})
	if err != nil {
		return fmt.Errorf("upload document: %w", docErr)
	}

	return ui.PrintDocumentSummary(app.Stdout, action, *(responseDoc.(*adminv1.DocumentResponse)))
}

func resolveWriteContent(cmd *cobra.Command, configuredEditor, slug string, args []string) ([]byte, error) {
	if len(args) == 2 {
		return input.Read(cmd.InOrStdin(), args[1])
	}

	if !isTTY(cmd.InOrStdin()) {
		return io.ReadAll(cmd.InOrStdin())
	}

	editorCommand := clieditor.Resolve(configuredEditor)
	content, changed, err := clieditor.EditBuffer(editorCommand, slug, nil)
	if err != nil {
		return nil, err
	}
	if !changed && len(strings.TrimSpace(string(content))) == 0 {
		return nil, fmt.Errorf("no content written")
	}

	return content, nil
}
