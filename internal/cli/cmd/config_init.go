// Copyright (c) 2026 Harsh Sharma <harsh@codelif.in>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// SPDX-License-Identifier: MIT

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

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	cliconfig "github.com/codelif/hostbin/internal/cli/config"
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
