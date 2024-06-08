package Simple

import (
	"fmt"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
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
			tab := p.GetTab()

			if len(args) == 0 {
				// Display help of the program.

				if p.Brief == "" {
					p.Println("Program:", p.Name)
				} else {
					p.Println("Program:", p.Name, "-", p.Brief)
				}
				p.Println()

				p.Println("Usage: ", p.Name, "(command) [arguments]")
				p.Println()

				if len(p.Description) > 0 {
					p.Println("Description:")
					for _, line := range p.Description {
						p.Println(tab, line)
					}
					p.Println()
				}

				table := make([][]string, 0, len(p.commands))
				for _, command := range p.commands {
					table = append(table, []string{command.Usage, command.Brief})
				}

				table, err := fs.TabAlign(table, 0, p.GetTabSize())
				if err != nil {
					return ue.NewErrWhile("tab aligning", err)
				}

				p.Println("Commands:")
				for _, row := range table {
					p.Println(tab, row[0], row[1])
				}
				p.Println()
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

				p.Println("Command:", command.Name)
				p.Println()

				p.Println("Usage: ", command.Usage)
				p.Println()

				if len(command.Description) > 0 {
					p.Println("Description:")
					for _, line := range command.Description {
						p.Println(tab, line)
					}
					p.Println()
				}
			}

			return nil
		},
	}
}
