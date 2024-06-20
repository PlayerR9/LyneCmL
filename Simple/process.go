package Simple

import (
	"fmt"

	ue "github.com/PlayerR9/MyGoLib/Units/errors"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// PhaseRunFunc is the function that runs the phase.
//
// Parameters:
//   - p: The program to run.
//   - prevRes: The result of the previous phase.
//   - data: The data to pass to the phase.
//
// Returns:
//   - any: The result of the phase.
//   - any: The data to pass to the next phase.
//   - error: The error that occurred while running the phase.
type PhaseRunFunc func(p *Program, prevRes any, data any) (any, any, error)

// PhaseSetupFunc is the function that sets up the program before running the phases.
//
// Parameters:
//   - p: The program to run.
//   - args: The arguments to pass to the program.
//   - data: The data to pass to the phases.
//
// Returns:
//   - any: The result of the setup.
//   - any: The data to pass to the phases.
//   - error: The error that occurred while setting up the program.
type PhaseSetupFunc func(p *Program, args []string, data any) (any, any, error)

// Phase is a phase of the program.
type Phase struct {
	// Name is the name of the phase.
	Name string

	// Short is a quick description of the phase.
	Short string

	// RunFunc is the function that runs the phase.
	RunFunc PhaseRunFunc
}

// MakeExecPhases creates a function that runs the phases.
//
// Parameters:
//   - p: The program to run.
//
// Returns:
//   - error: The error that occurred while running the phases.
func MakeExecPhases(setupFunc PhaseSetupFunc, phases ...*Phase) RunFunc {
	phases = us.FilterNilValues(phases)
	if len(phases) == 0 {
		if setupFunc == nil {
			return func(p *Program, args []string, data any) (any, error) {
				return nil, fmt.Errorf("no phases to run")
			}
		} else {
			return func(p *Program, args []string, data any) (any, error) {
				_, _, err := setupFunc(p, args, data)
				if err != nil {
					return nil, ue.NewErrWhile("setup", err)
				}

				return nil, nil
			}
		}
	}

	totalPhases := len(phases)

	if setupFunc == nil {
		return func(p *Program, args []string, data any) (any, error) {
			var res any = args
			var err error

			for i, phase := range phases {
				p.Printf("\nPhase (%d/%d): %s...\n\n", i+1, totalPhases, phase.Short)

				res, data, err = phase.RunFunc(p, res, data)
				if err != nil {
					return nil, ue.NewErrWhile(phase.Name, err)
				}
			}

			return nil, nil
		}
	} else {
		return func(p *Program, args []string, data any) (any, error) {
			var res any
			var err error

			res, data, err = setupFunc(p, args, data)
			if err != nil {
				return nil, ue.NewErrWhile("setup", err)
			}

			for i, phase := range phases {
				p.Printf("Phase (%d/%d): %s...", i+1, totalPhases, phase.Short)

				res, data, err = phase.RunFunc(p, res, data)
				if err != nil {
					return nil, ue.NewErrWhile(phase.Name, err)
				}
			}

			return nil, nil
		}
	}
}
