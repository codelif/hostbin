package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/spf13/cobra"

	cliapp "hostbin/internal/cli/app"
	cliconfig "hostbin/internal/cli/config"
)

func newConfigInitCommand(app *cliapp.App) *cobra.Command {
	var (
		serverURL    string
		authKey      string
		readKeyStdin bool
		force        bool
		skipCheck    bool
	)

	command := &cobra.Command{
		Use:   "init",
		Short: "Create a local hbcli config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.Store()
			if err != nil {
				return err
			}

			if _, err := store.Load(); err == nil && !force {
				return fmt.Errorf("config already exists at %s (use --force to overwrite)", store.Path())
			} else if err != nil && !errors.Is(err, cliconfig.ErrNotFound) && !force {
				return err
			}

			cfg := cliconfig.DefaultFile()
			cfg.ServerURL = strings.TrimSpace(serverURL)
			cfg.AuthKey = strings.TrimSpace(authKey)

			if readKeyStdin {
				value, err := readAllTrimmed(cmd.InOrStdin())
				if err != nil {
					return err
				}
				cfg.AuthKey = value
			}

			if cfg.ServerURL == "" {
				cfg.ServerURL = cliconfig.DefaultServerURL
			}

			if cfg.AuthKey == "" {
				if !isInteractive(cmd.InOrStdin(), app.Stdout) {
					return fmt.Errorf("auth_key is required; use --auth-key, --auth-key-stdin, or run interactively")
				}
				promptedServerURL, promptedAuthKey, err := app.Prompter.Init(cfg.ServerURL)
				if err != nil {
					return err
				}
				cfg.ServerURL = promptedServerURL
				cfg.AuthKey = promptedAuthKey
			}

			cfg = cfg.Normalized()
			if err := cfg.Validate(); err != nil {
				return err
			}

			if !skipCheck {
				if err := checkConfig(cmd.Context(), app, cfg, true); err != nil {
					return err
				}
			}

			if err := store.Save(cfg); err != nil {
				return err
			}

			_, _ = fmt.Fprintf(app.Stdout, "Wrote configuration\n  Path: %s\n", store.Path())

			if skipCheck {
				return nil
			}

			return nil
		},
	}

	command.Flags().StringVar(&serverURL, "server-url", "", "Admin server base URL")
	command.Flags().StringVar(&authKey, "auth-key", "", "Preshared auth key")
	command.Flags().BoolVar(&readKeyStdin, "auth-key-stdin", false, "Read auth key from stdin")
	command.Flags().BoolVar(&force, "force", false, "Overwrite any existing config file")
	command.Flags().BoolVar(&skipCheck, "skip-check", false, "Skip server and auth verification after writing config")

	return command
}

func isInteractive(stdin io.Reader, stdout io.Writer) bool {
	inFile, inOK := stdin.(*os.File)
	outFile, outOK := stdout.(*os.File)
	if !inOK || !outOK {
		return false
	}

	return term.IsTerminal(int(inFile.Fd())) && term.IsTerminal(int(outFile.Fd()))
}

func readAllTrimmed(reader io.Reader) (string, error) {
	buffer := bufio.NewReader(reader)
	value, err := io.ReadAll(buffer)
	if err != nil {
		return "", err
	}

	trimmed := strings.TrimSpace(string(value))
	if trimmed == "" {
		return "", fmt.Errorf("received empty value")
	}

	return trimmed, nil
}
