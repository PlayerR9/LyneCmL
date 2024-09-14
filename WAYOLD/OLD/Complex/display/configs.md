package display

import "github.com/gdamore/tcell"

var (
	// DefaultConfigs are the default configurations for the display.
	DefaultConfigs *Configs
)

func init() {
	DefaultConfigs = &Configs{
		TabSize: 3,
		Spacing: 1,
		Background: tcell.StyleDefault.
			Background(tcell.ColorGhostWhite).
			Foreground(tcell.ColorBlack),
	}
}

// Configs are the configurations for the display.
type Configs struct {
	// TabSize is the size of a tab character.
	TabSize int

	// Spacing is the spacing between columns.
	Spacing int

	// Background is the background color of the display.
	Background tcell.Style
}
