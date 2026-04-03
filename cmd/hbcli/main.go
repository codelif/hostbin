package main

import (
	"context"
	"fmt"
	"os"

	"github.com/codelif/hostbin/internal/cli/cmd"
)

func main() {
	root := cmd.NewRootCommand(os.Stdout, os.Stderr)
	if err := root.ExecuteContext(context.Background()); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
