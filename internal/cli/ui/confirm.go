package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func Confirm(r io.Reader, w io.Writer, prompt string) (bool, error) {
	if _, err := fmt.Fprintf(w, "%s [y/N]: ", prompt); err != nil {
		return false, err
	}

	line, err := bufio.NewReader(r).ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}

	value := strings.ToLower(strings.TrimSpace(line))
	return value == "y" || value == "yes", nil
}
