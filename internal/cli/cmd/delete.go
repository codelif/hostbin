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

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	"github.com/codelif/hostbin/internal/cli/ui"
)

func newDeleteCommand(app *cliapp.App) *cobra.Command {
	var yes bool

	command := &cobra.Command{
		Use:   "delete <slug>",
		Short: "Delete a document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			slug := args[0]
			if err := validateSlug(slug); err != nil {
				return err
			}

			_, _, client, err := loadClient(app)
			if err != nil {
				return err
			}

			if !yes {
				if !isTTY(cmd.InOrStdin()) {
					return fmt.Errorf("delete requires --yes in non-interactive mode")
				}
				confirmed, err := ui.Confirm(cmd.InOrStdin(), app.Stderr, "Delete "+slug+"?")
				if err != nil {
					return err
				}
				if !confirmed {
					_, _ = fmt.Fprintln(app.Stdout, "Aborted")
					return nil
				}
			}

			err = ui.RunSpinner(app.Stderr, "Deleting "+slug, func() error {
				_, err := client.DeleteDocument(cmd.Context(), slug)
				return err
			})
			if err != nil {
				return fmt.Errorf("delete document: %w", err)
			}

			return ui.PrintDeleteSummary(app.Stdout, slug)
		},
	}

	command.Flags().BoolVar(&yes, "yes", false, "Skip delete confirmation")
	return command
}
