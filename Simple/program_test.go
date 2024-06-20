package Simple

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestMakeProgram tests the MakeProgram function.
func TestMakeProgram(t *testing.T) {
	TestFlag := &Flag[string]{
		Name:        "test",
		Brief:       "A test flag.",
		Usage:       "",
		Description: NewDescription("This is a test flag."),
		DefaultVal:  "",
		ParseFunc: func(arg string) (string, error) {
			return arg, nil
		},
	}

	NowCmd := &Command{
		Name:   "now",
		Usages: nil,
		Brief:  "Prints the current date and time.",
		Argument: ExactlyNArgs(1).SetParseFunc(func(p *Program, args []string) (any, int, error) {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, 0, err
			}

			return num, 1, nil
		}),
		Run: func(p *Program, args []string, data any) (any, error) {
			str := strings.Repeat(" ", data.(int))

			err := p.Printf("The current date and time is:\n%s- %s\n",
				str,
				time.Now().Format(time.RFC1123),
			)
			if err != nil {
				return nil, err
			}

			err = p.SavePartial("now.txt")
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	NewYearCmd := &Command{
		Name:     "new-year",
		Usages:   nil,
		Brief:    "Prints the date and time of the next new year.",
		Argument: NoArgument,
		Run: func(p *Program, args []string, data any) (any, error) {
			now := time.Now()
			year := now.Year()

			x := TestFlag.GetValue()
			if x != "" {
				err := p.Println("Test flag value:", x)
				if err != nil {
					return nil, err
				}
			}

			newYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local)

			err := p.Println("The date and time of the next new year is:",
				newYear.Format(time.RFC1123),
			)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	NewYearCmd.SetFlags(
		TestFlag,
	)

	Program := &Program{
		Name:        "Test",
		Brief:       "A test program.",
		Description: NewDescription("This is a test program."),
		Version:     "v0.1.4",
	}

	Program.SetCommands(
		NowCmd,
		NewYearCmd,
	)

	err := ExecuteProgram(Program, []string{"Test", "now", "7"})
	if err != nil {
		t.Errorf("ExecuteProgram failed: %s", err.Error())
	}
}
