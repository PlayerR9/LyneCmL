package cml

import (
	"io"
	"strings"
	"text/tabwriter"
)

// TabWriter is a TabWriter.
type TabWriter struct {
	// w is the TabWriter.
	w *tabwriter.Writer
}

// Flush flushes the TabWriter.
//
// Returns:
//   - error: Any error that may have occurred.
func (tw TabWriter) Flush() error {
	if tw.w == nil {
		return nil
	}

	return tw.w.Flush()
}

// PrintRow prints a row of data.
//
// Parameters:
//   - columns: The columns to print.
//
// Returns:
//   - error: Any error that may have occurred.
func (tw TabWriter) PrintRow(columns ...string) error {
	data := []byte(strings.Join(columns, "\t") + "\n")

	n, err := tw.w.Write(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return io.ErrShortWrite
	}

	return nil
}
