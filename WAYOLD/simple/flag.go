package simple

import (
	"errors"
	"strings"
)

// Flag is a struct that represents a flag.
type Flag struct {
	// Name is the name of the flag.
	Name string

	// ShortName is the short name of the flag.
	ShortName rune

	// LongName is the long name of the flag.
	LongName string

	// value is the value of the flag.
	value bool
}

func (f Flag) Fix() error {
	if f.Name == "" {
		return errors.New("flag name cannot be empty")
	}

	if f.ShortName == '-' {
		return errors.New("short name cannot be '-'")
	}

	if strings.HasPrefix(f.LongName, "-") {
		return errors.New("long name cannot start with '-'")
	}

	if f.LongName == "" && f.ShortName == 0 {
		return errors.New("either long name or short name must be set")
	}

	return nil
}
