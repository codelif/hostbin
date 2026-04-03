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
	"os"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
)

func newGetCommand(app *cliapp.App) *cobra.Command {
	var savePath string

	command := &cobra.Command{
		Use:   "get <slug>",
		Short: "Print or save raw document content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateSlug(args[0]); err != nil {
				return err
			}

			_, _, client, err := loadClient(app)
			if err != nil {
				return err
			}

			content, err := client.GetDocumentContent(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get document content: %w", err)
			}

			if savePath == "" {
				_, err = app.Stdout.Write(content)
				return err
			}

			if err := os.WriteFile(savePath, content, 0o644); err != nil {
				return err
			}

			_, err = fmt.Fprintf(app.Stdout, "Saved to %s\n", savePath)
			return err
		},
	}

	command.Flags().StringVar(&savePath, "save", "", "Save content to a file instead of stdout")
	return command
}
