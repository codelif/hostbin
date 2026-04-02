package ui

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

var spinnerFrames = []string{"|", "/", "-", "\\"}

func RunSpinner(w io.Writer, label string, fn func() error) error {
	if !isTerminal(w) {
		_, _ = fmt.Fprintf(w, "%s...\n", label)
		return fn()
	}

	stop := make(chan struct{})
	var once sync.Once
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		index := 0
		for {
			select {
			case <-stop:
				_, _ = fmt.Fprint(w, "\r\033[2K")
				return
			case <-ticker.C:
				_, _ = fmt.Fprintf(w, "\r%s %s", spinnerFrames[index%len(spinnerFrames)], label)
				index++
			}
		}
	}()

	err := fn()
	once.Do(func() { close(stop) })
	return err
}

func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}
