package editor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func EditBuffer(editorCommand, slug string, initial []byte) ([]byte, bool, error) {
	tempFile, err := os.CreateTemp("", tempPattern(slug))
	if err != nil {
		return nil, false, err
	}
	tempPath := tempFile.Name()
	defer func() { _ = os.Remove(tempPath) }()

	if _, err := tempFile.Write(initial); err != nil {
		_ = tempFile.Close()
		return nil, false, err
	}
	if err := tempFile.Close(); err != nil {
		return nil, false, err
	}

	if err := Open(editorCommand, tempPath); err != nil {
		return nil, false, err
	}

	updated, err := os.ReadFile(tempPath)
	if err != nil {
		return nil, false, err
	}

	return updated, !bytes.Equal(initial, updated), nil
}

func tempPattern(slug string) string {
	clean := strings.ReplaceAll(filepath.Base(slug), string(filepath.Separator), "-")
	if clean == "." || clean == "" {
		clean = "document"
	}
	return fmt.Sprintf("hbcli-%s-*.txt", clean)
}
