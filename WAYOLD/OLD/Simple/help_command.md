package Simple

import (
	"strings"

	cut "github.com/PlayerR9/LyneCml/WAYOLD/OLD/Common/util"
	// ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
)

var (
	// HelpCmd is the help command.
	HelpCmd *Command
)

func init() {
	HelpCmd = &Command{
		Name: cut.HelpCmdOpcode,
		Usages: []string{
			"help [command]",
		},
		Brief: "Displays help information about the program or a specific command",
		Description: []string{
			strings.Join([]string{
				"If no command is specified, the help command will display help",
				"information about the program. Otherwise, the help command will",
				"display help information about the specified command.",
			}, " "),
		},
		Argument: AtMostNArgs(1),
		Run: func(p *Program, data any) error {
			// configs := p.Configs

			/* printer, trav := ffs.NewStdPrinter(
				ffs.NewFormatter(
					ffs.NewIndentConfig(configs.GetTabStr(), 0),
					ffs.NewFormatterConfig(configs.TabSize, configs.Spacing),
				),
			)
			*/
			args := data.([]string)

			if len(args) == 0 {
				// err := p.FString(trav)
				// if err != nil {
				// 	return err
				// }
			} else {
				// Display help of a specific command.
				// name := args[0]

				/* command, ok := p.commands[name]
				if !ok {
					return gcers.NewErrInvalidUsage(
						fmt.Errorf("command %q is not a valid command", name),
						"Use command \"help\" to see the list of available commands",
					)
				} */

				/* err := command.FString(trav, WithSpacing(configs.Spacing))
				if err != nil {
					return err
				} */
			}

			// pages := ffs.Stringfy(printer.GetPages(), 1)

			// err := p.Println(strings.Join(pages, "\f"))
			// if err != nil {
			// 	return err
			// }

			return nil
		},
	}

	err := HelpCmd.Fix()
	if err != nil {
		panic(err)
	}
}
