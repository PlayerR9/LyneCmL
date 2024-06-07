package Simple

import "fmt"

var (
	NoArgument *Argument
)

func init() {
	NoArgument = &Argument{bounds: [2]int{0, 0}}
}

type Argument struct {
	bounds [2]int
}

func AtLeastNArgs(n int) *Argument {
	if n < 0 {
		n = 0
	}

	return &Argument{bounds: [2]int{n, -1}}
}

func AtMostNArgs(n int) *Argument {
	if n <= 0 {
		return NoArgument
	} else {
		return &Argument{bounds: [2]int{0, n}}
	}
}

func ExactlyNArgs(n int) *Argument {
	if n <= 0 {
		return NoArgument
	} else {
		return &Argument{bounds: [2]int{n, n}}
	}
}

func RangeArgs(min, max int) *Argument {
	if min < 0 {
		min = 0
	}

	if max < 0 {
		max = 0
	}

	if min > max {
		min, max = max, min
	}

	if min == 0 && max == 0 {
		return NoArgument
	} else {
		return &Argument{bounds: [2]int{min, max}}
	}
}

func (a *Argument) validate(args []string) ([]string, error) {
	left := a.bounds[0]

	var right int

	if a.bounds[1] == -1 {
		right = len(args)
	} else {
		right = a.bounds[1]
	}

	if len(args) < left {
		return nil, fmt.Errorf("expected at least %d arguments, got %d", left, len(args))
	}

	if len(args) > right {
		return nil, fmt.Errorf("expected at most %d arguments, got %d", right, len(args))
	}

	return args, nil
}
