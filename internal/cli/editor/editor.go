package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Resolve(configured string) string {
	if strings.TrimSpace(configured) != "" {
		return configured
	}
	if visual := strings.TrimSpace(os.Getenv("VISUAL")); visual != "" {
		return visual
	}
	if editor := strings.TrimSpace(os.Getenv("EDITOR")); editor != "" {
		return editor
	}
	if runtime.GOOS == "windows" {
		return "notepad"
	}
	return "vi"
}

func Open(editorCommand, filePath string) error {
	parts := strings.Fields(editorCommand)
	if len(parts) == 0 {
		return fmt.Errorf("editor command is empty")
	}

	cmd := exec.Command(parts[0], append(parts[1:], filePath)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
