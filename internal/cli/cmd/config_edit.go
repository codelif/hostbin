package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
	cliconfig "hostbin/internal/cli/config"
	clieditor "hostbin/internal/cli/editor"
)

func newConfigEditCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Open the config file in an editor",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.Store()
			if err != nil {
				return err
			}

			cfg, err := store.Load()
			if err != nil {
				if err == cliconfig.ErrNotFound {
					return fmt.Errorf("no configuration found; run `hbcli config init` first")
				}
				return err
			}

			if err := clieditor.Open(clieditor.Resolve(cfg.Editor), store.Path()); err != nil {
				return err
			}

			_, err = store.Load()
			return err
		},
	}
}
