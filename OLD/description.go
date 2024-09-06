package LyneCmL

import (
	"strings"
)

// DescBuilder is a builder for a description.
type DescBuilder struct {
	// lines is the list of lines in the description.
	lines []string
}

// NewDescription creates a new description.
//
// Parameters:
//   - lines: The lines of the description.
//
// Returns:
//   - *Description: The new description.
func NewDescription(lines ...string) *DescBuilder {
	return &DescBuilder{
		lines: lines,
	}
}

// AddLine adds a line to the description.
//
// Parameters:
//   - sections: The sections of the line.
//
// Returns:
//   - *Description: The description.
func (d *DescBuilder) AddLine(sections ...string) *DescBuilder {
	d.lines = append(d.lines, strings.Join(sections, " "))

	return d
}

// Build builds the description.
//
// Returns:
//   - []string: The description.
func (d *DescBuilder) Build() []string {
	linesCopy := make([]string, len(d.lines))
	copy(linesCopy, d.lines)

	d.Reset()

	return linesCopy
}

// Reset resets the description.
func (d *DescBuilder) Reset() {
	for i := range d.lines {
		d.lines[i] = ""
	}

	d.lines = d.lines[:0]
}
