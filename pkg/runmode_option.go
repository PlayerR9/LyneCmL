package pkg

var (
	// DefaultOptions are the default options for a program.
	DefaultOptions *ProgramOptions
)

func init() {
	DefaultOptions = &ProgramOptions{
		TabSize: 3,
	}
}

// ProgramOptions are the optional options for a program.
type ProgramOptions struct {
	// TabSize is the size of a tab character.
	TabSize int
}

// fix fixes the program options.
func (po *ProgramOptions) fix() {
	if po.TabSize <= 0 {
		po.TabSize = 3
	}
}
