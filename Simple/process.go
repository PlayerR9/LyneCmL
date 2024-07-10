package Simple

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

var (
	// DefaultPhaseRunFunc is the default function that runs the phase.
	DefaultPhaseRunFunc PhaseRunFunc
)

func init() {
	DefaultPhaseRunFunc = func(_ *Program, prevRes any, data any) (any, error) {
		return prevRes, nil
	}
}

// PhaseRunFunc is the function that runs the phase.
//
// Parameters:
//   - p: The program to run.
//   - prevRes: The result of the previous phase.
//   - data: The data to pass to the phase.
//
// Returns:
//   - any: The result of the phase to pass to the next phase.
//   - error: The error that occurred while running the phase.
type PhaseRunFunc func(p *Program, prevRes any, data any) (any, error)

// PhaseSetupFunc is the function that sets up the program before running the phases.
//
// Parameters:
//   - p: The program to run.
//   - data: The data to pass to the phases.
//
// Returns:
//   - any: The result of the setup.
//   - error: The error that occurred while setting up the program.
type PhaseSetupFunc func(p *Program, data any) (any, error)

// Phase is a phase of the program.
type Phase struct {
	// Name is the name of the phase. Only used for error messages.
	Name string

	// Short is a quick description of the phase. Used for printing out the phase.
	//
	// Format:
	//   "Phase (<i>/<N>): <Short>..."
	//
	// Where:
	//   - <i> is the index of the phase. (i.e., the order in which the phases are run)
	//   - <N> is the total number of phases.
	//   - <Short> is the brief description of the phase. This field is used here.
	//
	// Example:
	//
	// If <Short> is "Taking input", then the output will be: "Phase (1/3): Taking input..."
	Short string

	// RunFunc is the function that runs the phase.
	//
	// Nil run functions are not filtered out. However, they will be converted to the
	// default DefaultPhaseRunFunc which simply forwards the previous result and data
	// to the next phase.
	RunFunc PhaseRunFunc
}

// Fix implements the common.Fixer interface.
//
// This never errors.
//
// Behaviors:
//   - The Name and Short fields are trimmed of whitespace on both ends.
//   - If the RunFunc is nil, then it is set to the DefaultPhaseRunFunc.
func (p *Phase) Fix() error {
	p.Name = strings.TrimSpace(p.Name)
	p.Short = strings.TrimSpace(p.Short)

	if p.RunFunc == nil {
		p.RunFunc = DefaultPhaseRunFunc
	}

	return nil
}

// PhaseString creates a string for the phase.
//
// Format:
//
//	"Phase (<i>/<N>): <Brief>..."
//
// Where:
//   - <i> is the index of the phase. (i.e., the order in which the phases are run)
//   - <N> is the total number of phases.
func PhaseString(i, total int, brief string) string {
	var builder strings.Builder

	builder.WriteRune('\n')
	builder.WriteString("Phase (")
	builder.WriteString(strconv.Itoa(i + 1))
	builder.WriteRune('/')
	builder.WriteString(strconv.Itoa(total))
	builder.WriteString("): ")
	builder.WriteString(brief)
	builder.WriteString("...")
	builder.WriteRune('\n')

	return builder.String()
}

// MakeExecPhases creates a function that runs the phases.
//
//	Args -> Setup -> Phase1 -> Phase2 -> ... -> PhaseN
//
// Where:
//   - Args: The arguments passed to the program.
//   - Setup: The setup function that sets up the program before running the phases.
//   - Phase1, Phase2, ..., PhaseN: The phases to run.
//
// Each phase prints out the phase name before running the phase according to the
// following format:
//
//	Phase (i/N): Name...
//
// Here, i is the index of the phase, N is the total number of phases, and
// Name is the brief description of the phase.
//
// Like so: "Phase (1/3): Taking input..."
//
// Parameters:
//   - setupFunc: The function that sets up the program before running the phases.
//   - phases: The phases to run.
//
// Returns:
//   - RunFunc: The function that runs the phases.
//
// Behaviors:
//   - The phases are fixed before running.
//   - If there are no phases, then the setup function is run. If
//     there is no setup function, then the RunFunc will do nothing.
//   - If there are phases, then the setup function is run before the
//     phases are run in order as they are passed.
//   - The res of a phase is passed to the next phase as the prevRes. Also,
//     if the setup function is provided, then the res of the setup function
//     is passed to the first phase as the prevRes. Otherwise, the args are
//     passed to the first phase as the prevRes.
func MakeExecPhases(setupFunc PhaseSetupFunc, phases ...*Phase) RunFunc {
	for _, phase := range phases {
		phase.Fix()
	}

	if len(phases) == 0 {
		if setupFunc == nil {
			return func(p *Program, data any) error {
				return nil
			}
		} else {
			return func(p *Program, data any) error {
				_, err := setupFunc(p, data)
				if err != nil {
					return uc.NewErrWhile("setup", err)
				}

				return nil
			}
		}
	}

	totalPhases := len(phases)

	if setupFunc == nil {
		return func(p *Program, data any) error {
			var res any = data

			for i, phase := range phases {
				str := PhaseString(i+1, totalPhases, phase.Short)

				err := p.Println(str)
				if err != nil {
					return uc.NewErrWhile("print", err)
				}

				res, err = phase.RunFunc(p, res, data)
				if err != nil {
					return uc.NewErrWhile(phase.Name, err)
				}
			}

			return nil
		}
	} else {
		return func(p *Program, data any) error {
			res, err := setupFunc(p, data)
			if err != nil {
				return uc.NewErrWhile("setup", err)
			}

			for i, phase := range phases {
				str := PhaseString(i+1, totalPhases, phase.Short)

				err := p.Println(str)
				if err != nil {
					return uc.NewErrWhile("print", err)
				}

				res, err = phase.RunFunc(p, res, data)
				if err != nil {
					return uc.NewErrWhile(phase.Name, err)
				}
			}

			return nil
		}
	}
}

// ExecProcess is a process that runs a command.
type ExecProcess struct {
	// args are the arguments to pass to the command.
	args []string

	// data is the data to pass to the command.
	data any

	// cmd is the command to run.
	cmd *Command
}

// NewExecProcess creates a new ExecProcess.
//
// Parameters:
//   - args: The arguments to pass to the command.
//   - data: The data to pass to the command.
//   - cmd: The command to run.
//
// Returns:
//   - *ExecProcess: The new ExecProcess. Nil if the command is nil.
func NewExecProcess(args []string, data any, cmd *Command, flagLeft []*Flag) *ExecProcess {
	if cmd == nil {
		return nil
	}

	for _, flag := range flagLeft {
		if flag.Argument == nil {
			flag.value = false
		} else {
			flag.value = flag.Argument.defaultVal
		}
	}

	return &ExecProcess{
		args: args,
		data: data,
		cmd:  cmd,
	}
}

// Execute runs the command.
//
// Parameters:
//   - p: The program to run the command on.
//
// Returns:
//   - error: An error if the command failed to run.
func (ep *ExecProcess) Execute(p *Program) error {
	if p == nil {
		return uc.NewErrNilParameter("program")
	}

	ok := p.display.IsDone()
	if ok {
		return errors.New("program is done")
	}

	var builder strings.Builder

	builder.WriteString("Running command ")
	qn := strconv.Quote(ep.cmd.Name)
	builder.WriteString(qn)
	builder.WriteString(" ...")

	err := p.Println(builder.String())
	if err != nil {
		return fmt.Errorf("error printing: %w", err)
	}

	err = p.Println()
	if err != nil {
		return fmt.Errorf("error printing: %w", err)
	}

	p.flags = ep.cmd.flags

	err = ep.cmd.Run(p, ep.data)
	if err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
