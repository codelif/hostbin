package cmd

import (
	"io"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
)

func NewRootCommand(stdout, stderr io.Writer) *cobra.Command {
	application := cliapp.New(stdout, stderr)

	root := &cobra.Command{
		Use:           "hbcli",
		Short:         "Manage hostbin configuration and documents",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVar(&application.ConfigPath, "config", "", "Path to hbcli config file")
	root.AddCommand(
		newConfigCommand(application),
		newListCommand(application),
		newInfoCommand(application),
		newGetCommand(application),
		newNewCommand(application),
		newPutCommand(application),
		newDeleteCommand(application),
		newEditCommand(application),
	)

	return root
}
