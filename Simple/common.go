package Simple

import (
	com "github.com/PlayerR9/LyneCmL/Simple/common"
	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

///////////////////////////////////////////////////////

const (
	// ConfigLoc is the default location of the configuration files.
	ConfigLoc string = "configs"
)

// CmlComponent is a component of a CML program.
type CmlComponent interface {
	// GenerateUsage generates the usage of the component.
	//
	// Returns:
	//   - []string: The usage of the component.
	GenerateUsage() []string

	ffs.FStringer

	com.Fixer
}

const (
	// HelpCommandName is the name of the help command.
	HelpCmdOpcode string = "help"
)

var (
	glTableAligner *com.TableAligner
)

// SetTableAligner sets the table aligner.
//
// No operation will be performed if the table aligner is already set.
// If the table aligner is nil, then no operation will be performed.
//
// Parameters:
//   - ta: The table aligner to set.
func SetTableAligner(ta *com.TableAligner) {
	if glTableAligner != nil {
		return
	}
	if ta == nil {
		return
	}

	glTableAligner = ta
}
