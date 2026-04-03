package cmd

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/term"

	cliapp "github.com/codelif/hostbin/internal/cli/app"
	cliconfig "github.com/codelif/hostbin/internal/cli/config"
	clihttp "github.com/codelif/hostbin/internal/cli/http"
	"github.com/codelif/hostbin/internal/domain/slugs"
)

func loadClient(app *cliapp.App) (*cliconfig.Store, cliconfig.File, *clihttp.Client, error) {
	store, err := app.Store()
	if err != nil {
		return nil, cliconfig.File{}, nil, err
	}

	cfg, err := store.Load()
	if err != nil {
		if err == cliconfig.ErrNotFound {
			return nil, cliconfig.File{}, nil, fmt.Errorf("no configuration found; run `hbcli config init`")
		}
		return nil, cliconfig.File{}, nil, err
	}

	client, err := app.Client(cfg)
	if err != nil {
		return nil, cliconfig.File{}, nil, err
	}

	return store, cfg, client, nil
}

func validateSlug(value string) error {
	if err := slugs.Validate(value, nil); err != nil {
		return fmt.Errorf("invalid slug %q", value)
	}

	return nil
}

func isTTY(reader io.Reader) bool {
	file, ok := reader.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}
