package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	cliconfig "github.com/codelif/hostbin/internal/cli/config"
)

func newConfigCheckCommand(app *cliapp.App) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Verify server connectivity and authentication",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := app.Store()
			if err != nil {
				return err
			}

			return runConfigCheck(cmd.Context(), app, store, false)
		},
	}
}

func runConfigCheck(ctx context.Context, app *cliapp.App, store *cliconfig.Store, afterSave bool) error {
	cfg, err := store.Load()
	if err != nil {
		if err == cliconfig.ErrNotFound {
			return fmt.Errorf("no configuration found; run `hbcli config init`")
		}
		return err
	}

	return checkConfig(ctx, app, cfg, afterSave)
}

func checkConfig(ctx context.Context, app *cliapp.App, cfg cliconfig.File, announce bool) error {
	if ctx == nil {
		ctx = context.Background()
	}
	timeout, err := cfg.Duration()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := app.Client(cfg)
	if err != nil {
		return err
	}

	if announce {
		_, _ = fmt.Fprintln(app.Stdout, "Checking configuration...")
	}
	if _, err := client.Health(ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	_, _ = fmt.Fprintln(app.Stdout, "OK server reachable")

	if _, err := client.AuthCheck(ctx); err != nil {
		return fmt.Errorf("auth check failed: %w", err)
	}
	_, _ = fmt.Fprintln(app.Stdout, "OK authentication valid")

	return nil
}
