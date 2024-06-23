package configs

import (
	"path"

	uo "github.com/PlayerR9/MyGoLib/Utility/object"
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

// Configer is a configuration.
type Configer interface {
	// Default sets the configuration to the default values.
	Default()

	uo.Fixer
}
