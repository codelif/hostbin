package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
	clieditor "hostbin/internal/cli/editor"
	"hostbin/internal/cli/ui"
	"hostbin/internal/protocol/adminv1"
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
