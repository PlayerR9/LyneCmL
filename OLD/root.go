package LyneCmL

import (
	"fmt"
	"log"
	"os"

	prs "github.com/PlayerR9/LyneCml/OLD/Common/Parser"
	cms "github.com/PlayerR9/LyneCml/OLD/Simple"
	pd "github.com/PlayerR9/LyneCml/OLD/Simple/display"
	gcers "github.com/PlayerR9/go-commons/errors"
)

// Fix is a function that fixes the program. This should be
// always called before running the program as it will fix
// any errors in the program.
//
// Parameters:
//   - p: The program to fix.
//
// Returns:
//   - error: An error if the program failed to fix.
func Fix(p *cms.Program) error {
	if p == nil {
		return gcers.NewErrNilParameter("program")
	}

	err := p.Fix()
	if err != nil {
		return fmt.Errorf("invalid program: %w", err)
	}

	return nil
}

// ExecuteProgram runs the program.
//
// Parameters:
//   - p: The program to run.
//   - args: The arguments to run the program with. (i.e., os.Args)
//
// Returns:
//   - error: An error if the program failed to run.
//
// Behaviors:
//   - If no program is given, it will do nothing and return nil.
//
// Example:
//
//	var Program *Simple.Program
//
//	func init() {
//		Program = &Simple.Program{
//			Name:  "Program",
//			Brief: "A program that can be run.",
//			Description: []string{
//				"Program is a program that can be run.",
//				"It can run commands and display information.",
//			},
//			Version: "0.0.1",
//			Options: nil,
//		}
//
//		Program.AddCommand(
//			// your commands here
//		)
//
//		Program.AddConfig(
//			// your configs here
//		)
//
//		err := Fix(Program)
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	func main() {
//		err := ExecuteProgram(Program, os.Args)
//		DefaultExitSequence(err) // or any other exit sequence
//	}
func ExecuteProgram(p *cms.Program, args []string) error {
	defer p.SaveConfigs() // Save the configs after running the program.

	if p.Name == "" {
		p.Name = args[0]
	}

	args = args[1:]

	cmdMap := p.GetCmdMap()

	list, err := prs.ParseArgs(cmdMap, args)
	if err != nil {
		return fmt.Errorf("in parsing arguments: %w", err)
	}

	display := pd.NewDisplay(p.Configs, log.New(os.Stdout, "["+p.Name+"]: ", log.LstdFlags))

	p.SetDisplay(display)

	display.Start()
	defer display.Close()

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		pe := gcers.NewErrPanic(r)

		err := p.Panic(pe)
		if err != nil {
			panic(err)
		}
	}()

	for len(list) > 0 {
		top := list[len(list)-1]
		list = list[:len(list)-1]

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
