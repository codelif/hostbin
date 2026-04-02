package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
	cliconfig "hostbin/internal/cli/config"
)

func newConfigSetCommand(app *cliapp.App) *cobra.Command {
	var readKeyStdin bool

	command := &cobra.Command{
		Use:   "set <key> [value]",
		Short: "Write a config value",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.Store()
			if err != nil {
				return err
			}

			cfg, err := store.Load()
			if err != nil {
				if err == cliconfig.ErrNotFound {
					cfg = cliconfig.File{}
				} else {
					return err
				}
			}

			key := args[0]
			value := ""
			if len(args) == 2 {
				value = args[1]
			}

			if readKeyStdin {
				if key != "auth_key" {
					return fmt.Errorf("--auth-key-stdin can only be used with auth_key")
				}
				value, err = readAllTrimmed(cmd.InOrStdin())
				if err != nil {
					return err
				}
			}

			if strings.TrimSpace(value) == "" {
				current, err := cfg.Get(key)
				if err != nil {
					return err
				}
				if !isTerminalReader(cmd.InOrStdin()) || !isTerminalWriter(app.Stdout) {
					return fmt.Errorf("value is required for %s in non-interactive mode", key)
				}
				value, err = app.Prompter.Value(key, current)
				if err != nil {
					return err
				}
			}

			if err := cfg.Set(key, value); err != nil {
				return err
			}

			if err := store.Save(cfg); err != nil {
				return err
			}

			_, err = fmt.Fprintf(app.Stdout, "Updated %s\n", key)
			return err
		},
	}

	command.Flags().BoolVar(&readKeyStdin, "auth-key-stdin", false, "Read auth_key from stdin")
	return command
}

func isTerminalReader(reader io.Reader) bool {
	file, ok := reader.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}

func isTerminalWriter(writer io.Writer) bool {
	file, ok := writer.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}
