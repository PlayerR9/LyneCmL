package style

import "github.com/gdamore/tcell"

//go:generate stringer -type ColorType

type ColorType int

const (
	NormalText ColorType = iota
	SuccessText
	ErrorText
	WarningText
	DebugText
	InfoText
)

var (
	// LightModeStyle is a default style for light mode.
	LightModeStyle Style[ColorType]
)

func init() {
	LightModeStyle = NewStyle[ColorType]()

	LightModeStyle.SetNormalText(tcell.StyleDefault.Background(tcell.ColorGhostWhite).Foreground(tcell.ColorBlack))
	LightModeStyle.SetSuccessText(tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack).Bold(true))
	LightModeStyle.SetErrorText(tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite).Bold(true))
	LightModeStyle.AddStyle(WarningText, tcell.StyleDefault.Background(tcell.ColorOrange).Foreground(tcell.ColorBlack).Bold(true))
	LightModeStyle.AddStyle(DebugText, tcell.StyleDefault.Background(tcell.ColorGhostWhite).Foreground(tcell.ColorDarkGrey))
	LightModeStyle.AddStyle(InfoText, tcell.StyleDefault.Background(tcell.ColorGhostWhite).Foreground(tcell.ColorDarkGoldenrod))
}

var (
	// DarkModeStyle is a default style for dark mode.
	DarkModeStyle Style[ColorType]
)

func init() {
	DarkModeStyle = NewStyle[ColorType]()

	DarkModeStyle.SetNormalText(tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhiteSmoke))
	DarkModeStyle.SetSuccessText(tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorLightGreen).Bold(true))
	DarkModeStyle.SetErrorText(tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.ColorLightCoral).Bold(true))
	DarkModeStyle.AddStyle(WarningText, tcell.StyleDefault.Background(tcell.ColorDarkOrange).Foreground(tcell.ColorLightSalmon).Bold(true))
	DarkModeStyle.AddStyle(DebugText, tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorLightGrey))
	DarkModeStyle.AddStyle(InfoText, tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorLightYellow))
}

var (
	// GreenStyle is a default style for green text.
	GreenStyle Style[ColorType]
)

func init() {
	GreenStyle = NewStyle[ColorType]()

	GreenStyle.SetNormalText(tcell.StyleDefault.Background(tcell.ColorSpringGreen).Foreground(tcell.ColorDarkGrey))
	GreenStyle.SetSuccessText(tcell.StyleDefault.Background(tcell.ColorLime).Foreground(tcell.ColorDarkGrey).Bold(true))
	GreenStyle.SetErrorText(tcell.StyleDefault.Background(tcell.ColorSpringGreen).Foreground(tcell.ColorRed).Bold(true))
	GreenStyle.AddStyle(WarningText, tcell.StyleDefault.Background(tcell.ColorSpringGreen).Foreground(tcell.ColorOrange).Bold(true))
	GreenStyle.AddStyle(DebugText, tcell.StyleDefault.Background(tcell.ColorSpringGreen).Foreground(tcell.ColorGrey))
	GreenStyle.AddStyle(InfoText, tcell.StyleDefault.Background(tcell.ColorSpringGreen).Foreground(tcell.ColorYellow))
}
