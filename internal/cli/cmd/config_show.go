package cmd

import (
	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	"github.com/codelif/hostbin/internal/cli/ui"
)

func newConfigShowCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Print the current config summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.Store()
			if err != nil {
				return err
			}
			cfg, err := store.Load()
			if err != nil {
				return err
			}

			return ui.PrintConfigSummary(app.Stdout, store.Path(), cfg)
		},
	}
}
