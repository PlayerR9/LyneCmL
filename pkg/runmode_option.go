package pkg

var (
	// DefaultOptions are the default options for a program.
	//
	// === Default Options ===
	// TabSize: 3
	// Spacing: 1
	DefaultOptions *Configurations
)

func init() {
	DefaultOptions = &Configurations{
		TabSize: 3,
		Spacing: 1,
	}
}

// Configurations are the optional options for a program.
type Configurations struct {
	// TabSize is the size of a tab character.
	TabSize int

	// Spacing is the spacing between columns.
	Spacing int
}

// fix fixes the program options.
func (po *Configurations) fix() {
	if po.TabSize <= 0 {
		po.TabSize = 3
	}

	if po.Spacing <= 0 {
		po.Spacing = 1
	}
}
