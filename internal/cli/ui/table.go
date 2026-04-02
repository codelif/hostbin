package ui

import (
	"fmt"
	"io"
	"text/tabwriter"

	"hostbin/internal/cli/format"
	"hostbin/internal/protocol/adminv1"
)

func PrintDocumentTable(w io.Writer, documents []adminv1.DocumentResponse) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "SLUG\tSIZE\tUPDATED"); err != nil {
		return err
	}
	for _, doc := range documents {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\n", doc.Slug, format.Bytes(doc.SizeBytes), format.Timestamp(doc.UpdatedAt)); err != nil {
			return err
		}
	}

	return tw.Flush()
}
