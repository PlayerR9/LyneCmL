package LyneCmL

import (
	"strconv"
	"strings"
	"testing"
	"time"

	cms "github.com/PlayerR9/LyneCmL/Simple"
	com "github.com/PlayerR9/LyneCmL/Simple/common"
)

// TestMakeProgram tests the MakeProgram function.
func TestMakeProgram(t *testing.T) {
	TestFlag := &cms.Flag[string]{
		LongName:    "test",
		Brief:       "A test flag.",
		Usage:       "",
		Description: com.NewDescription("This is a test flag.").Build(),
		DefaultVal:  "",
		ParseFunc: func(arg string) (string, error) {
			return arg, nil
		},
	}

	NowCmd := &cms.Command{
		Name:   "now",
		Usages: nil,
		Brief:  "Prints the current date and time.",
		Argument: cms.ExactlyNArgs(1).SetParseFunc(func(args []string) (any, error) {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, err
			}

			return num, nil
		}),
		Run: func(p *cms.Program, data any) error {
			str := strings.Repeat(" ", data.(int))

			err := p.Printf("The current date and time is:\n%s- %s\n",
				str,
				time.Now().Format(time.RFC1123),
			)
			if err != nil {
				return err
			}

			err = p.SavePartial("now.txt")
			if err != nil {
				return err
			}

			return nil
		},
	}

	NewYearCmd := &cms.Command{
		Name:     "new-year",
		Usages:   nil,
		Brief:    "Prints the date and time of the next new year.",
		Argument: cms.NoArgument,
		Run: func(p *cms.Program, data any) error {
			now := time.Now()
			year := now.Year()

			x := TestFlag.Value
			if x != "" {
				err := p.Println("Test flag value:", x)
				if err != nil {
					return err
				}
			}

			newYear := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.Local)

			err := p.Println("The date and time of the next new year is:",
				newYear.Format(time.RFC1123),
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	NewYearCmd.SetFlags(
		TestFlag,
	)

	Program := &cms.Program{
		Name:        "Test",
		Brief:       "A test program.",
		Description: com.NewDescription("This is a test program.").Build(),
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
