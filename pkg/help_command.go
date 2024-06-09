package pkg

import (
	"fmt"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
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
		Description: NewDescription(
			"If no command is specified, the help command will display help information about the program.",
			"Otherwise, the help command will display help information about the specified command.",
		),
		Argument: AtMostNArgs(1),
		Run: func(p *Program, args []string) error {
			printer := ffs.NewStdPrinter(
				ffs.NewFormatter(
					ffs.NewIndentConfig(p.GetTab(), 0),
				),
			)

			if len(args) == 0 {
				err := ffs.Apply(printer, p)
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

				err := ffs.Apply(printer, command)
				if err != nil {
					return err
				}
			}

			pages := ffs.Stringfy(printer.GetPages())

			p.Println(strings.Join(pages, "\f"))

			return nil
		},
	}
}
