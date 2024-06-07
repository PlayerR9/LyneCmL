package Simple

var (
	HelpCommandCmd *Command
)

func init() {
	HelpCommandCmd = &Command{
		Name: "help",
		Description: []string{
			"Displays help information about the program",
			"Usage: ZesseConv help",
		},
		Argument: NoArgument,
		Run: func(p *Program, args []string) error {
			p.Println("Program:", p.Name)
			p.Println()

			for _, line := range p.Description {
				p.Println(line)
			}
			p.Println()

			p.Println("Usage:")
			p.Println("  ", p.Name, "(command) [arguments]")
			p.Println()

			p.Println("Commands:")

			for _, command := range p.commands {
				p.Println(command.Name)

				for _, line := range command.Description {
					p.Println("  ", line)
				}
				p.Println()
			}

			return nil
		},
	}
}
