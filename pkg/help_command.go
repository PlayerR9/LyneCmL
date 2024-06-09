package pkg

import (
	"fmt"

	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

var (
	// HelpCommandCmd is the help command.
	HelpCommandCmd *Command
)

func init() {
	HelpCommandCmd = &Command{
		Name:  HelpCommandOpcode,
		Usage: "help [command]",
		Brief: "Displays help information about the program or a specific command",
		Description: []string{
			"If no command is specified, the help command will display help information about the program.",
			"Otherwise, the help command will display help information about the specified command.",
		},
		Argument: AtMostNArgs(1),
		Run: func(p *Program, args []string) error {
			var lines []string

			if len(args) == 0 {
				var err error

				lines, err = p.DisplayHelp()
				if err != nil {
					return err
				}
			} else {
				// Display help of a specific command.
				name := args[0]

				command, ok := p.commands[name]
				if !ok {
					return ue.NewErrInvalidUsage(
						fmt.Errorf("command %q is not a valid command", name),
						"Use command \"help\" to see the list of available commands",
					)
				}

				lines = command.DisplayHelp()
			}

			for _, line := range lines {
				p.Println(line)
			}
			p.Println()

			return nil
		},
	}
}
