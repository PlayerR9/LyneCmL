package configs

import (
	"strings"
)

///////////////////////////////////////////////////////

// Config are the configurations for a program.
type DisplayConfigs struct {
	// TabSize is the size of a tab character.
	TabSize int `json:"tab_size"`

	// Spacing is the spacing between columns.
	Spacing int `json:"spacing"`
}

// Fix implements Configer interface.
//
// This never returns an error.
func (dc *DisplayConfigs) Fix() error {
	if dc.TabSize <= 0 {
		dc.TabSize = 3
	}

	if dc.Spacing <= 0 {
		dc.Spacing = 1
	}

	return nil
}

// Default implements Configer interface.
func (dc *DisplayConfigs) Default() {
	dc.Spacing = 1
	dc.TabSize = 3
}

// GetSpacingStr gets the spacing string.
//
// Returns:
//   - string: The spacing string.
func (dc *DisplayConfigs) GetSpacingStr() string {
	str := strings.Repeat(" ", dc.Spacing)

	return str
}

// GetTabStr gets the tab string.
//
// Returns:
//   - string: The tab string.
func (dc *DisplayConfigs) GetTabStr() string {
	spacing := strings.Repeat(" ", dc.Spacing)

	str := strings.Repeat(spacing, dc.TabSize)

	return str
}
