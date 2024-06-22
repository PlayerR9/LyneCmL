package configs

import (
	"encoding/json"
	"strings"
)

///////////////////////////////////////////////////////

// Config are the configurations for a program.
type DisplayConfigs struct {
	// TabSize is the size of a tab character.
	TabSize int

	// Spacing is the spacing between columns.
	Spacing int
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
func (dc *DisplayConfigs) Default() Configer {
	config := &DisplayConfigs{
		TabSize: 3,
		Spacing: 1,
	}

	return config
}

// MarshalJSON implements Configer interface.
func (dc *DisplayConfigs) MarshalJSON() ([]byte, error) {
	type Alias struct {
		TabSize int `json:"tab_size"`
		Spacing int `json:"spacing"`
	}

	a := &Alias{
		TabSize: dc.TabSize,
		Spacing: dc.Spacing,
	}

	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UnmarshalJSON implements Configer interface.
func (dc *DisplayConfigs) UnmarshalJSON(data []byte) error {
	type Alias struct {
		TabSize int `json:"tab_size"`
		Spacing int `json:"spacing"`
	}

	a := Alias{
		TabSize: dc.TabSize,
		Spacing: dc.Spacing,
	}

	err := json.Unmarshal(data, &a)
	if err != nil {
		return err
	}

	dc.TabSize = a.TabSize
	dc.Spacing = a.Spacing

	return nil
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
