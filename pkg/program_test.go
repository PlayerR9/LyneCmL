package pkg

import (
	"testing"
	"time"
)

// TestMakeProgram tests the MakeProgram function.
func TestMakeProgram(t *testing.T) {
	NowCmd := &Command{
		Name:     "now",
		Usage:    "now",
		Brief:    "Prints the current date and time.",
		Argument: NoArgument,
		Run: func(p *Program, args []string, data any) (any, error) {
			p.Println("The current date and time is:", time.Now().Format(time.RFC1123))

			return nil, nil
		},
	}

	NewYearCmd := &Command{
		Name:     "new-year",
		Usage:    "new-year",
		Brief:    "Prints the date and time of the next new year.",
		Argument: NoArgument,
		Run: func(p *Program, args []string, data any) (any, error) {
			now := time.Now()
			year := now.Year()

			newYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local)

			p.Println("The date and time of the next new year is:", newYear.Format(time.RFC1123))

			return nil, nil
		},
	}

	Program := GenerateProgram(
		&Command{},
		"v0.1.3",
		nil,
		NowCmd,
		NewYearCmd,
	)

	ExecuteProgram(Program, []string{"Test", "now"})
}
