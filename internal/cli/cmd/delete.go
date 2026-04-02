package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
	"hostbin/internal/cli/ui"
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
