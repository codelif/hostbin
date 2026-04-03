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
