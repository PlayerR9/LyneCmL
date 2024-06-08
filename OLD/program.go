package OLD

import (
	"log"
	"os"
	"strings"

	sdb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	ue "github.com/PlayerR9/MyGoLib/Units/errors"
)

// Memorer is the memory of the program.
type Memorer interface{}

// ProgramMode is the mode of the program.
//
// Parameters:
//   - p: The program to apply the mode to.
//
// Returns:
//   - error: The error that occurred while applying the mode.
type ProgramMode func(*Program) error

// WithCustomLogger creates a program mode that uses a custom logger.
//
// Parameters:
//   - logger: The custom logger to use.
//
// Returns:
//   - ProgramMode: The program mode that uses the custom logger.
func WithCustomLogger(logger *log.Logger) ProgramMode {
	return func(p *Program) error {
		// debug := p.logger.GetDebugMode()

		if logger == nil {
			logger = p.makeDefaultLogger()
		}

		p.logger = sdb.NewDebugger(logger)

		// p.logger.ToggleDebugMode(debug)

		return nil
	}
}

// WithActiveLogger creates a program mode that uses an active logger.
//
// Returns:
//   - ProgramMode: The program mode that uses the active logger.
func WithActiveLogger() ProgramMode {
	return func(p *Program) error {
		p.logger.ToggleDebugMode(true)
		return nil
	}
}

// Program is a program that can be run.
type Program struct {
	// Name is the name of the program.
	Name string

	// Init is the pre-run function that initializes the program.
	// It should be used for configuration reading and setting up the program.
	Init func(*Program) (Memorer, error)

	// Memory is the memory of the program (i.e., the configuration of the program).
	Memory Memorer

	// IsPedantic is the flag that determines whether the program should be pedantic.
	// If true, the program won't ignore ErrIgnorable errors.
	IsPedantic bool

	// stdio is the channel for the standard input/output.
	stdio *sdb.Debugger

	// logger is the channel for the log messages.
	logger *sdb.Debugger

	// commandList is the list of commands that the program can run.
	commandList map[string]*Command
}

// makeDefaultLogger is a helper function that creates a default logger.
//
// Returns:
//   - *log.Logger: The default logger.
func (p *Program) makeDefaultLogger() *log.Logger {
	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(p.Name)
	builder.WriteString("]: ")

	logger := log.New(os.Stdout, builder.String(), log.LstdFlags)

	return logger
}

// programSetup sets up the program.
//
// Parameters:
//   - opts: The options to apply to the program.
//
// Returns:
//   - error: The error that occurred while setting up the program.
func (p *Program) programSetup(opts []ProgramMode) error {
	p.stdio = sdb.NewDebugger(nil)
	p.stdio.ToggleDebugMode(true)

	p.logger = sdb.NewDebugger(p.makeDefaultLogger())

	for i, opt := range opts {
		err := opt(p)
		if err != nil {
			return ue.NewErrWhileAt("applying", i+1, "option", err)
		}
	}

	p.stdio.Start()
	defer p.stdio.Close()

	p.logger.Start()
	defer p.logger.Close()

	memory, reason := p.Init(p)
	err := pedanticErrorEval(p.IsPedantic, reason, NewErrFailedInitialization(reason))
	if err != nil {
		return err
	}

	p.Memory = memory

	return nil
}

// Println prints a line to the standard output.
//
// Parameters:
//   - a: The items to print.
func (p *Program) Println(a ...interface{}) {
	p.stdio.Println(a...)
}

// Printf prints a formatted line to the standard output.
//
// Parameters:
//   - format: The format of the line.
//   - a: The items to print.
func (p *Program) Printf(format string, a ...interface{}) {
	p.stdio.Printf(format, a...)
}

// Logln logs a line.
//
// Parameters:
//   - a: The items to log.
func (p *Program) Logln(a ...interface{}) {
	p.logger.Println(a...)
}

// Logf logs a formatted line.
//
// Parameters:
//   - format: The format of the line.
//   - a: The items to log.
func (p *Program) Logf(format string, a ...interface{}) {
	p.logger.Printf(format, a...)
}

// GetCommand gets a command by its opcode.
//
// Parameters:
//   - opcode: The opcode of the command.
//
// Returns:
//   - *Command: The command with the opcode.
func (p *Program) GetCommand(opcode string) *Command {
	cmd, ok := p.commandList[opcode]
	if !ok {
		return nil
	}

	return cmd
}

// AddCommand adds a command to the program.
//
// Parameters:
//   - cmd: The command to add.
func (p *Program) AddCommand(cmd *Command) {
	if p.commandList == nil {
		p.commandList = make(map[string]*Command)
	}

	name := cmd.Name

	p.commandList[name] = cmd
}

// Run runs the program.
//
// Parameters:
//   - p: The program to run.
//   - opts: The options to apply to the program.
//
// Returns:
//   - error: The error that occurred while running the program.
func Run(p *Program, opts ...ProgramMode) error {
	// 1. Setup the program.
	err := p.programSetup(opts)
	if err != nil {
		return err
	}

	// 2. Argument parsing.
	if len(os.Args) < 2 {
		return NewErrNoCommandProvided()
	}

	opcode := os.Args[1]

	cmd := p.GetCommand(opcode)
	if cmd == nil {
		return NewErrCommandNotFound(opcode)
	}

	// 3. Command execution.
	reason := cmd.Execute(p, os.Args[2:])

	err = pedanticErrorEval(p.IsPedantic, reason, NewErrCommandFailed(opcode, reason))
	if err != nil {
		return err
	}

	return nil
}
