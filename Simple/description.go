package Simple

import (
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

// Description is a description of a command or flag.
type Description struct {
	// lines is the list of lines in the description.
	lines []string
}

// FString formats the description using the given traversor.
//
// Parameters:
//   - trav: The traversor to use to format the description.
//   - opts: The options to use to format the description.
//
// Returns:
//   - error: An error if the description failed to format.
func (d *Description) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	err := trav.AddLines(d.lines)
	if err != nil {
		return err
	}

	return nil
}

// NewDescription creates a new description.
//
// Parameters:
//   - lines: The lines of the description.
//
// Returns:
//   - *Description: The new description.
func NewDescription(lines ...string) *Description {
	return &Description{
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
func (d *Description) AddLine(sections ...string) *Description {
	d.lines = append(d.lines, strings.Join(sections, " "))

	return d
}
