package cmd

import (
	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
)

func newConfigCommand(app *cliapp.App) *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "Manage hbcli configuration",
	}

	command.AddCommand(
		newConfigInitCommand(app),
		newConfigGetCommand(app),
		newConfigSetCommand(app),
		newConfigEditCommand(app),
		newConfigPathCommand(app),
		newConfigShowCommand(app),
		newConfigCheckCommand(app),
	)

	return command
}
