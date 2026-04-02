package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
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
