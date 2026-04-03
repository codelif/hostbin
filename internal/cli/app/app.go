package app

import (
	"io"

	cliconfig "github.com/codelif/hostbin/internal/cli/config"
	clihttp "github.com/codelif/hostbin/internal/cli/http"
)

type App struct {
	Stdout io.Writer
	Stderr io.Writer

	ConfigPath string
	Prompter   cliconfig.Prompter
}

func New(stdout, stderr io.Writer) *App {
	return &App{
		Stdout:   stdout,
		Stderr:   stderr,
		Prompter: cliconfig.HuhPrompter{},
	}
}

func (a *App) Store() (*cliconfig.Store, error) {
	return cliconfig.NewStore(a.ConfigPath)
}

func (a *App) Client(cfg cliconfig.File) (*clihttp.Client, error) {
	return clihttp.New(cfg)
}
