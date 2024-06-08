package pkg

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// parseArgs is a helper function that parses the arguments of a program.
//
// Parameters:
//   - p: The program to parse the arguments for.
//   - args: The arguments to parse.
//
// Returns:
//   - *Command: The command that was parsed.
//   - []string: The parsed arguments.
//   - error: The error that occurred while parsing the arguments.
func parseArgs(p *Program, args []string) (*Command, []string, error) {
	if len(args) < 1 {
		return nil, nil, errors.New("no command provided")
	}

	command, ok := p.commands[args[0]]
	if !ok {
		return nil, nil, fmt.Errorf("unknown command: %q", args[0])
	}

	parsedArgs, err := command.Argument.validate(args[1:])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid arguments: %w", err)
	}

	err = command.Argument.parseFunc(p, parsedArgs)
	if err != nil {
		return nil, nil, err
	}

	return command, parsedArgs, nil
}

// runBody is a helper function that runs the body of a program.
//
// Parameters:
//   - p: The program to run.
//   - args: The arguments to run the program with.
//   - opts: The options to apply to the program.
//
// Returns:
//   - error: The error that occurred while running the program.
func runBody(p *Program, args []string) error {
	command, parsedArgs, err := parseArgs(p, args[1:])
	if err != nil {
		return ue.NewErrInvalidUsage(
			err,
			"Use \"help\" to see a list of available commands.",
		)
	}

	p.buffer = sfb.NewBuffer[string]()

	p.buffer.Start()
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			str, ok := p.buffer.Receive()
			if !ok {
				break
			}

			if strings.HasSuffix(str, "\n") {
				fmt.Print(str)
			} else {
				fmt.Println(str)
			}
		}
	}()

	err = command.Run(p, parsedArgs)

	p.buffer.Close()
	wg.Wait()

	if err != nil {
		return ue.NewErrWhile("running command", err)
	}

	return nil
}

// Run runs the program.
//
// Parameters:
//   - args: The arguments to run the program with. (i.e., os.Args)
func (p *Program) Run(args []string) {
	p.AddCommands(HelpCommandCmd)

	p.fix(args[0])

	err := runBody(p, args)
	if err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("Program finished successfully")
	}
	fmt.Println()

	fmt.Println("Press ENTER to exit...")
	fmt.Scanln()
}
