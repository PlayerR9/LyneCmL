package util

import (
	// ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	gcf "github.com/PlayerR9/go-commons/fixer"
)

// CmlComponent is a component of a CML program.
type CmlComponent interface {
	// GenerateUsage generates the usage of the component.
	//
	// Returns:
	//   - []string: The usage of the component.
	GenerateUsage() []string

	// ffs.FStringer

	gcf.Fixer
}

const (
	// HelpCommandName is the name of the help command.
	HelpCmdOpcode string = "help"
)
