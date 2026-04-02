package ui

import (
	"fmt"
	"io"

	cliconfig "hostbin/internal/cli/config"
	"hostbin/internal/cli/format"
	"hostbin/internal/protocol/adminv1"
)

func PrintConfigSummary(w io.Writer, path string, cfg cliconfig.File) error {
	normalized := cfg.Normalized()
	_, err := fmt.Fprintf(w,
		"hbcli configuration\n  path:       %s\n  server_url: %s\n  auth_key:   %s\n  timeout:    %s\n  editor:     %s\n  color:      %s\n",
		path,
		valueOrUnset(normalized.ServerURL),
		valueOrUnset(cliconfig.MaskSecret(normalized.AuthKey)),
		normalized.Timeout,
		valueOrUnset(normalized.Editor),
		normalized.Color,
	)
	return err
}

func valueOrUnset(value string) string {
	if value == "" {
		return "(unset)"
	}

	return value
}

func PrintDocumentSummary(w io.Writer, action string, doc adminv1.DocumentResponse) error {
	_, err := fmt.Fprintf(w,
		"%s %s\n  URL: %s\n  Size: %s\n  Updated: %s\n",
		action,
		doc.Slug,
		doc.URL,
		format.Bytes(doc.SizeBytes),
		format.Timestamp(doc.UpdatedAt),
	)
	return err
}

func PrintDocumentInfo(w io.Writer, doc adminv1.DocumentResponse) error {
	_, err := fmt.Fprintf(w,
		"%s\n  URL: %s\n  Size: %s\n  SHA256: %s\n  Created: %s\n  Updated: %s\n",
		doc.Slug,
		doc.URL,
		format.Bytes(doc.SizeBytes),
		doc.SHA256,
		doc.CreatedAt,
		doc.UpdatedAt,
	)
	return err
}

func PrintDeleteSummary(w io.Writer, slug string) error {
	_, err := fmt.Fprintf(w, "Deleted %s\n", slug)
	return err
}
