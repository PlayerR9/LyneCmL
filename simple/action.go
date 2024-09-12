package simple

import (
	"fmt"

	"github.com/gdamore/tcell"
)

type Action interface {
}

type ActPrint struct {
	args  []any
	style tcell.Style
}

func NewActPrint(style tcell.Style, args ...any) *ActPrint {
	return &ActPrint{
		args:  args,
		style: style,
	}
}

type ActPrintf struct {
	format string
	args   []any
	style  tcell.Style
}

func NewActPrintf(style tcell.Style, format string, args ...any) *ActPrintf {
	return &ActPrintf{
		style:  style,
		format: format + "\n",
		args:   args,
	}
}

func do_action(c *Context, action Action) error {
	switch act := action.(type) {
	case *ActPrint:
		str := fmt.Sprint(act.args...)

		for x, char := range []rune(str) {
			c.screen.DrawCell(x, 0, char, act.style)
		}

		// TODO: Finish this
	case *ActPrintf:
		str := fmt.Sprintf(act.format, act.args...)

		for x, char := range []rune(str) {
			c.screen.DrawCell(x, 1, char, act.style)
		}

		// TODO: Finish this
	default:
		return fmt.Errorf("invalid action type: %T", action)
	}

	return nil
}
