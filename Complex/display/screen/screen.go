package screen

import (
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell"

	ddt "github.com/PlayerR9/MyGoLib/Display/drawtable"
	rws "github.com/PlayerR9/MyGoLib/Safe/RWSafe"
)

// Display represents a display that can draw elements to the screen.
type Display struct {
	// screen is the screen of the display.
	screen tcell.Screen

	// width is the width of the display.
	width int

	// height is the height of the display.
	height int

	// wg is the wait group of the display.
	wg sync.WaitGroup

	// table is the draw table of the display.
	table *ddt.DrawTable

	// element is the element to draw.
	element *rws.Subject[ddt.Displayer]

	// errChan is the channel of errors.
	errChan chan error

	// keyChan is the channel of key events.
	keyChan chan tcell.EventKey

	// shouldClose is the subject of whether the display should close.
	shouldClose *rws.Subject[bool]

	// bgStyle is the background style of the display.
	bgStyle tcell.Style
}

// NewDisplay creates a new display with the given background style.
//
// Parameters:
//   - bgStyle: The background style of the display.
//
// Returns:
//   - *Display: The new display.
//   - error: An error if the display could not be created.
func NewDisplay(bgStyle tcell.Style) (*Display, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.SetStyle(bgStyle)
	screen.Clear()

	width, height := screen.Size()

	table := ddt.NewDrawTable(width, height)

	return &Display{
		screen:  screen,
		width:   width,
		height:  height,
		table:   table,
		bgStyle: bgStyle,
	}, nil
}

// Start starts the display.
func (d *Display) Start() {
	d.errChan = make(chan error, 1)
	d.keyChan = make(chan tcell.EventKey)

	d.shouldClose = rws.NewSubject[bool](false)
	d.element = rws.NewSubject[ddt.Displayer](nil)
}

// Close closes the display.
func (d *Display) Close() {
	d.shouldClose.Set(true)

	d.wg.Wait()

	close(d.errChan) // Check this
	d.errChan = nil

	d.screen.Fini()

	close(d.keyChan)
	d.keyChan = nil
}

// ReceiveErr receives an error from the display.
//
// Returns:
//   - error: The error.
//   - bool: True if the error was received, false otherwise.
func (d *Display) ReceiveErr() (error, bool) {
	err, ok := <-d.errChan
	if !ok {
		return nil, false
	}

	return err, true
}

// Draw draws an element to the display.
//
// Parameters:
//   - elem: The element to draw.
func (d *Display) Draw(elem ddt.Displayer) {
	d.element.Set(elem)

	d.drawScreen()
}

/*
// mainListener is a helper method that listens for events.
func (d *Display) mainListener() {
	defer d.wg.Done()

	for {
		select {
		case <-time.After(time.Microsecond * 100):
			if d.shouldClose.Get() {
				return
			}
		case ev := <-d.evChan:

		}
	}
}
*/

// resizeEvent is a helper method that handles a resize event.
func (d *Display) resizeEvent() {
	d.width, d.height = d.screen.Size()

	d.table = ddt.NewDrawTable(d.width, d.height)
}

// drawScreen is a helper method that draws the screen.
func (d *Display) drawScreen() {
	d.screen.Clear()

	elem := d.element.Get()

	if elem != nil {
		xCoord := 2
		yCoord := 2

		err := elem.Draw(d.table, &xCoord, &yCoord)
		if err != nil {
			d.errChan <- fmt.Errorf("error drawing element: %w", err)
		}

		height := d.table.GetHeight()
		width := d.table.GetWidth()

		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				cell := d.table.GetAt(j, i)

				if cell == nil {
					continue
				}

				runes, err := cell.Runes(1, 1)
				if err != nil {
					d.errChan <- fmt.Errorf("error getting runes: %w", err)
				}

				d.screen.SetContent(j, i, runes[0][0], nil, cell.GetStyle())
			}
		}
	}

	d.screen.Show()
	time.Sleep(time.Millisecond * 100)
}
