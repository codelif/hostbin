package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	"github.com/codelif/hostbin/internal/cli/ui"
)

func newListCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List documents",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _, client, err := loadClient(app)
			if err != nil {
				return err
			}

			documents, err := client.ListDocuments(cmd.Context())
			if err != nil {
				return fmt.Errorf("list documents: %w", err)
			}

			return ui.PrintDocumentTable(app.Stdout, documents)
		},
	}
}
