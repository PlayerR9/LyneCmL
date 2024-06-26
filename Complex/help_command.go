package Complex

import (
	"fmt"
	"strings"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

var (
	// HelpCmd is the help command.
	HelpCmd *Command
)

func init() {
	HelpCmd = &Command{
		Name:  HelpCmdOpcode,
		Usage: "help [command]",
		Brief: "Displays help information about the program or a specific command",
		Description: NewDescription(
			"If no command is specified, the help command will display help information about the program.",
			"Otherwise, the help command will display help information about the specified command.",
		),
		Argument: AtMostNArgs(1),
		Run: func(p *Program, args []string, data any) (any, error) {
			printer, trav := ffs.NewStdPrinter(
				ffs.NewFormatter(
					ffs.NewIndentConfig(p.GetTab(), 0),
					ffs.NewFormatterConfig(p.GetTabSize(), p.GetSpacing()),
				),
			)

			if len(args) == 0 {
				err := p.FString(trav)
				if err != nil {
					return nil, err
				}
			} else {
				// Display help of a specific command.
				name := args[0]

				command, ok := p.commands[name]
				if !ok {
					return nil, uc.NewErrInvalidUsage(
						fmt.Errorf("command %q is not a valid command", name),
						"Use command \"help\" to see the list of available commands",
					)
				}

				err := command.FString(trav)
				if err != nil {
					return nil, err
				}
			}

			pages := ffs.Stringfy(printer.GetPages(), 1)

			err := p.Println(strings.Join(pages, "\f"))
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}
}
