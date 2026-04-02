package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
	"hostbin/internal/cli/ui"
)

func newInfoCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "info <slug>",
		Short: "Show document metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateSlug(args[0]); err != nil {
				return err
			}

			_, _, client, err := loadClient(app)
			if err != nil {
				return err
			}

			doc, err := client.GetDocument(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("get document info: %w", err)
			}

			return ui.PrintDocumentInfo(app.Stdout, *doc)
		},
	}
}
