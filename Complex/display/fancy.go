package display

import (
	"fmt"

	"github.com/gdamore/tcell"
)

type FancyDisplay struct {
	screen tcell.Screen

	width  int
	height int
}

func NewFancyDisplay() (*FancyDisplay, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	err = screen.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize screen: %w", err)
	}

	width, height := screen.Size()

	return &FancyDisplay{
		screen: screen,
		width:  width,
		height: height,
	}, nil
}
