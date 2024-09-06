package display

import (
	"path"
	"strings"
)

const (
	// DisplayConfig is the display configuration.
	DisplayConfig string = "display"

	// InternalsConfig is the internals configuration.
	InternalsConfig string = "internals"

	// ConfigDir is the directory of the configuration files.
	ConfigDir string = "configs"
)

var (
	// ConfigLoc is the location of the configuration files.
	ConfigLoc string
)

func init() {
	ConfigLoc = path.Join(ConfigDir, "config.json")
}

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

// NewDisplayConfigs creates a new display configuration.
//
// Returns:
//   - *DisplayConfigs: The new display configuration.
func NewDisplayConfigs() *DisplayConfigs {
	dc := &DisplayConfigs{
		TabSize: 3,
		Spacing: 1,
	}
	return dc
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
