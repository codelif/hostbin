package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	cliconfig "github.com/codelif/hostbin/internal/cli/config"
)

func newConfigGetCommand(app *cliapp.App) *cobra.Command {
	var raw bool

	command := &cobra.Command{
		Use:   "get <key>",
		Short: "Read a config value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.Store()
			if err != nil {
				return err
			}
			cfg, err := store.Load()
			if err != nil {
				return err
			}

			value, err := cfg.Get(args[0])
			if err != nil {
				return err
			}
			if args[0] == "auth_key" && !raw {
				value = cliconfig.MaskSecret(value)
			}

			_, err = fmt.Fprintln(app.Stdout, value)
			return err
		},
	}

	command.Flags().BoolVar(&raw, "raw", false, "Print secret values without masking")
	return command
}
