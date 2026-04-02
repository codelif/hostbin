package input

import (
	"fmt"
	"io"
	"os"
)

func Read(stdin io.Reader, filePath string) ([]byte, error) {
	if filePath == "-" {
		return io.ReadAll(stdin)
	}

	if filePath == "" {
		return nil, fmt.Errorf("input source is empty")
	}

	return os.ReadFile(filePath)
}
