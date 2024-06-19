package display

var (
	// DefaultConfigs are the default configurations for the display.
	DefaultConfigs *Configs
)

func init() {
	DefaultConfigs = &Configs{
		TabSize: 3,
		Spacing: 1,
	}
}

// Configs are the configurations for the display.
type Configs struct {
	// TabSize is the size of a tab character.
	TabSize int

	// Spacing is the spacing between columns.
	Spacing int
}
