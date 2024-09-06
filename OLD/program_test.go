package LyneCmL

import (
	"strconv"
	"strings"
	"testing"
	"time"

	cms "github.com/PlayerR9/LyneCml/OLD/Simple"
)

// TestMakeProgram tests the MakeProgram function.
func TestMakeProgram(t *testing.T) {
	TestFlag := &cms.Flag{
		LongName:    "test",
		Brief:       "A test flag.",
		Usages:      nil,
		Description: NewDescription("This is a test flag.").Build(),
		Argument:    cms.DefaultFlagArgument,
	}

	NowCmd := &cms.Command{
		Name:   "now",
		Usages: nil,
		Brief:  "Prints the current date and time.",
		Argument: cms.ExactlyNArgs(1).WithParseFunc(func(args []string) (any, error) {
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

			// err = p.SavePartial("now.txt")
			// if err != nil {
			// 	return err
			// }

			return nil
		},
	}

	NowCmd.SetFlags(
		TestFlag,
	)

	NewYearCmd := &cms.Command{
		Name:     "new-year",
		Usages:   nil,
		Brief:    "Prints the date and time of the next new year.",
		Argument: cms.NoArgument,
		Run: func(p *cms.Program, data any) error {
			now := time.Now()
			year := now.Year()

			x := TestFlag.Value().(string)
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
		Description: NewDescription("This is a test program.").Build(),
		Version:     "v0.1.10",
	}

	Program.SetCommands(
		NowCmd,
		NewYearCmd,
	)

	err := Fix(Program)
	if err != nil {
		t.Fatalf("Fix failed: %s", err.Error())
	}

	// err = ExecuteProgram(Program, []string{"Test", "now", "7", "--test", "yes"})
	// if err != nil {
	// 	t.Fatalf("ExecuteProgram failed: %s", err.Error())
	// }

	err = ExecuteProgram(Program, []string{"Test", "help"})
	if err != nil {
		t.Errorf("ExecuteProgram failed: %s", err.Error())
	}
}

func TestHelpCommand(t *testing.T) {
	Program := &cms.Program{
		Name:        "Test",
		Brief:       "A test program.",
		Description: NewDescription("This is a test program.").Build(),
		Version:     "v0.1.4",
	}

	err := Fix(Program)
	if err != nil {
		t.Fatalf("Fix failed: %s", err.Error())
	}

	// err = ExecuteProgram(Program, []string{"Test", "now", "7", "--test", "yes"})
	// if err != nil {
	// 	t.Fatalf("ExecuteProgram failed: %s", err.Error())
	// }

	err = ExecuteProgram(Program, []string{"Test", "help"})
	if err != nil {
		t.Errorf("ExecuteProgram failed: %s", err.Error())
	}

	t.Fatalf("Test failed")
}
