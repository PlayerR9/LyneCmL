package Complex

import (
	"fmt"
	"sync"

	util "github.com/PlayerR9/LyneCmL/Complex/util"
	llq "github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	"golang.org/x/net/context"
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

	p.buffer = sfb.NewBuffer[any]()

	p.history = llq.NewSafeQueue[string]()

	p.buffer.Start()

	ctx, cancel := context.WithCancel(context.Background())
	p.ctx = ctx

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, ok := p.buffer.Receive()
				if !ok {
					return
				}

				err := p.msgHandler(msg)
				if err != nil {
					fmt.Println("Error:", err)
					cancel()
					return
				}
			}
		}
	}()

	args = args[1:]

	cmd, ok := p.commands[args[0]]
	if !ok {
		return util.NewErrUnknownCommand(args[0])
	}

	_, err := parseArgs(p, args[1:], cmd)
	if err != nil {
		cancel()
		return fmt.Errorf("in command %q: %w", cmd.Name, err)
	}

	p.buffer.Close()
	wg.Wait()

	cancel()

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
