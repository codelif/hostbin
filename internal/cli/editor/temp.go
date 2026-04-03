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
