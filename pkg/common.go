package pkg

import (
	"fmt"
	"strings"
	"sync"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

const (
	// HelpCommandName is the name of the help command.
	HelpCmdOpcode string = "help"
)

// GenerateProgram creates a new program with the given command, version, options,
// and subcommands.
//
// Parameters:
//   - cmd: The main/privileged command of the program.
//   - version: The version of the program.
//   - opts: The options of the program.
//   - subCmds: The subcommands of the program.
//
// Returns:
//   - *Program: The new program.
//
// Behaviors:
//   - If cmd is nil, it will be set to a new command.
//   - nil subcommands will be filtered out.
//   - If a subcommand has the same name as another subcommand, the first one
//     will be kept.
func GenerateProgram(cmd *Command, version string, opts *Configurations, subCmds ...*Command) *Program {
	if cmd == nil {
		// Deal with nil command
	} else {
		cmd.fix()

		cmd.Usage = strings.Join([]string{
			"Usage:", cmd.Name, "(command)", "[arguments]",
		}, " ")

		cmd.Argument = ProgramDefaultArgument
		cmd.Run = ProgramDefaultRunFunc
	}

	if opts == nil {
		opts = DefaultOptions
	} else {
		opts.fix()
	}

	p := &Program{
		Command: cmd,
		Version: version,
		commands: map[string]*Command{
			HelpCmdOpcode: HelpCmd,
		},
		Options: opts,
	}

	subCmds = us.FilterNilValues(subCmds)
	for _, subCmd := range subCmds {
		_, ok := p.commands[subCmd.Name]
		if !ok {
			p.commands[subCmd.Name] = subCmd
		}
	}

	for _, command := range p.commands {
		command.fix()
		if command.Name == "" {
			continue
		}

		_, ok := p.commands[command.Name]
		if !ok {
			p.commands[command.Name] = command
		}
	}

	return p
}

// ExecuteProgram runs the program.
//
// Parameters:
//   - p: The program to run.
//   - args: The arguments to run the program with. (i.e., os.Args)
//
// Returns:
//   - error: An error if the program failed to run.
func ExecuteProgram(p *Program, args []string) error {
	if p.Name == "" {
		p.Name = args[0]
		p.Usage = strings.Join([]string{
			"Usage:", args[0], "(command)", "[arguments]",
		}, " ")
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

			space := strings.Repeat(" ", p.Options.Spacing)

			str, _ = fs.FixTabStop(0, p.Options.TabSize, space, str)

			if strings.HasSuffix(str, "\n") {
				fmt.Print(str)
			} else {
				fmt.Println(str)
			}

			wg.Add(1)

			// Enqueue the line to the history in a separate goroutine
			// to prevent blocking the buffer.
			go func(line string) {
				defer wg.Done()

				line = strings.TrimSuffix(line, "\n")

				p.history.Enqueue(line)
			}(str)
		}
	}()

	_, err := parseArgs(p, args[1:], p.Command)

	p.buffer.Close()
	wg.Wait()

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
