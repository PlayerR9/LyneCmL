package simple

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
//   - *DescBuilder: The new description builder.
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
//   - *DescBuilder: The description. Nil only if receiver is nil.
func (d *DescBuilder) AddLine(sections ...string) *DescBuilder {
	if d == nil {
		return nil
	}

	d.lines = append(d.lines, strings.Join(sections, " "))

	return d
}

// Build builds the description.
//
// Returns:
//   - []string: The description.
func (d DescBuilder) Build() []string {
	line_copy := make([]string, len(d.lines))
	copy(line_copy, d.lines)

	return line_copy
}

// Reset resets the description.
func (d *DescBuilder) Reset() {
	if d == nil {
		return
	}

	if len(d.lines) > 0 {
		for i := range d.lines {
			d.lines[i] = ""
		}

		d.lines = d.lines[:0]
	}
}
