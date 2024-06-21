package LyneCmL

import (
	"fmt"
	"log"
	"os"

	prs "github.com/PlayerR9/LyneCmL/Common/Parser"
	cms "github.com/PlayerR9/LyneCmL/Simple"
	com "github.com/PlayerR9/LyneCmL/Simple/common"
	cnf "github.com/PlayerR9/LyneCmL/Simple/configs"
	pd "github.com/PlayerR9/LyneCmL/Simple/display"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// ExecuteProgram runs the program.
//
// Parameters:
//   - p: The program to run.
//   - args: The arguments to run the program with. (i.e., os.Args)
//
// Returns:
//   - error: An error if the program failed to run.
func ExecuteProgram(p *cms.Program, args []string) error {
	p.Fix()

	if p.Name == "" {
		p.Name = args[0]
	}

	args = args[1:]

	cmdMap := p.GetCmdMap()

	list, err := prs.ParseArgs(cmdMap, args)
	if err != nil {
		return fmt.Errorf("in parsing arguments: %w", err)
	}

	p.Options = cnf.NewConfig("configs", 0644)

	dc := p.Options.GetConfigs(cnf.DisplayConfig).(*cnf.DisplayConfigs)

	cms.SetTableAligner(com.NewTableAligner(dc.TabSize))

	display := pd.NewDisplay(dc, log.New(os.Stdout, "["+p.Name+"]: ", log.LstdFlags))

	p.SetDisplay(display)

	display.Start()
	defer display.Close()

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		pe := ue.NewErrPanic(r)

		err := p.Panic(pe)
		if err != nil {
			panic(err)
		}
	}()

	for {
		top, ok := list.Pop()
		if !ok {
			break
		}

		err := top.Execute(p)
		if err != nil {
			return fmt.Errorf("error executing command: %w", err)
		}
	}

	return nil
}

// DefaultExitSequence is the default exit sequence for a program.
//
// Parameters:
//   - err: The error that occurred while running the program.
//
// Behaviors:
//   - If err is nil, it will print "Program ran successfully."
//   - If err is not nil, it will print "Error: <error>".
//   - It will print "Press ENTER to exit...".
//   - It will wait for the user to press ENTER.
func DefaultExitSequence(err error) {
	if err != nil {
		fmt.Println("Error:", err.Error())
		fmt.Println("Use \"help\" to see a list of available commands.")
	} else {
		fmt.Println("Program ran successfully.")
	}

	fmt.Println()

	fmt.Println("Press ENTER to exit...")
	fmt.Scanln()
}
