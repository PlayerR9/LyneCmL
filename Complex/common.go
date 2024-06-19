package Complex

import (
	"fmt"

	pd "github.com/PlayerR9/LyneCmL/Complex/display"
	util "github.com/PlayerR9/LyneCmL/Complex/util"
)

const (
	// HelpCommandName is the name of the help command.
	HelpCmdOpcode string = "help"
)

// ExecuteProgram runs the program.
//
// Parameters:
//   - p: The program to run.
//   - args: The arguments to run the program with. (i.e., os.Args)
//
// Returns:
//   - error: An error if the program failed to run.
func ExecuteProgram(p *Program, args []string) error {
	p.fix()

	if p.Name == "" {
		p.Name = args[0]
	}

	displayConfigs := &pd.Configs{
		TabSize: p.Options.TabSize,
		Spacing: p.Options.Spacing,
	}

	display, err := pd.NewDisplay(displayConfigs)
	if err != nil {
		return err
	}

	p.display = display

	p.display.Start()
	defer p.display.Close()

	args = args[1:]

	cmd, ok := p.commands[args[0]]
	if !ok {
		return util.NewErrUnknownCommand(args[0])
	}

	_, err = parseArgs(p, args[1:], cmd)
	if err != nil {
		return fmt.Errorf("in command %q: %w", cmd.Name, err)
	}

	return err
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
