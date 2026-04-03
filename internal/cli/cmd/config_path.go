package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
)

func newConfigPathCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the active config file path",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.Store()
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(app.Stdout, store.Path())
			return err
		},
	}
}
