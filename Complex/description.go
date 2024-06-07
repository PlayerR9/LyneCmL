package Complex

import "strings"

// Description is a description of a command or flag.
type Description struct {
	// lines is the list of lines in the description.
	lines []string
}

// NewDescription creates a new description.
//
// Returns:
//   - *Description: The new description.
func NewDescription() *Description {
	return &Description{
		lines: make([]string, 0),
	}
}

// AddLine adds a line to the description.
//
// Parameters:
//   - sections: The sections of the line.
//
// Returns:
//   - *Description: The description.
func (d *Description) AddLine(sections ...string) *Description {
	d.lines = append(d.lines, strings.Join(sections, " "))

	return d
}
