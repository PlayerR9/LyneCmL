package simple

import (
	"os"

	"github.com/PlayerR9/LyneCml/simple/internal"
	cmls "github.com/PlayerR9/LyneCml/style"
	ds "github.com/PlayerR9/display/screen"
	gcers "github.com/PlayerR9/go-commons/errors"
	"github.com/gdamore/tcell"
)

type ExitSequence struct {
	style cmls.Style[cmls.ColorType]
	err   error
}

func (es ExitSequence) Draw(screen ds.Drawable, xCoord, yCoord *int) error {
	if screen == nil {
		return nil
	}

	y := 0

	if es.err == nil {
		for x, c := range []rune("Program ran successfully.") {
			screen.DrawCell(x, y, c, es.style.SuccessText)
		}

		y++
	} else {
		for x, c := range []rune(es.err.Error()) {
			screen.DrawCell(x, y, c, es.style.ErrorText)
		}

		y++

		switch err := es.err.(type) {
		case *gcers.Err[internal.ErrorCode]:
			y++

			for x, c := range []rune("Suggestions:") {
				screen.DrawCell(x, y, c, es.style.NormalText)
			}

			for _, suggestion := range err.Suggestions {
				y++

				for x, c := range []rune(suggestion) {
					screen.DrawCell(x+3, y, c, es.style.NormalText)
				}
			}

			y++
		}
	}

	y++

	for x, c := range []rune("Press ENTER to exit...") {
		screen.DrawCell(x, y, c, es.style.NormalText)
	}

	return nil
}

// DefaultExitSequence is the default exit sequence for the program.
//
// Parameters:
//   - err: The error that occurred. If nil, the program will exit with code 0.
func DefaultExitSequence(p *Program, err error) {
	if p == nil {
		return
	}

	defer p.screen.Close()

	y := 0

	var exit_code int

	if err == nil {
		for x, c := range []rune("Program ran successfully.") {
			p.screen.SetCell(x, y, c, p.style.SuccessText)
		}

		y++
		exit_code = 0
	} else {
		for x, c := range []rune(err.Error()) {
			p.screen.SetCell(x, y, c, p.style.ErrorText)
		}

		y++

		switch err := err.(type) {
		case *gcers.Err[internal.ErrorCode]:
			y++

			for x, c := range []rune("Suggestions:") {
				p.screen.SetCell(x, y, c, p.style.NormalText)
			}

			for _, suggestion := range err.Suggestions {
				y++

				for x, c := range []rune(suggestion) {
					p.screen.SetCell(x+3, y, c, p.style.NormalText)
				}
			}

			y++

			exit_code = int(err.Code) + 2
		default:
			exit_code = 1
		}
	}

	y++

	for x, c := range []rune("Press ENTER to exit...") {
		p.screen.SetCell(x, y, c, p.style.NormalText)
	}

	for {
		key, ok := p.screen.ListenForKey()
		if !ok {
			break
		}

		if key.Key() == tcell.KeyEnter {
			break
		}
	}

	os.Exit(exit_code)
}
