package cml

var (
	// NoArguments is the default argument.
	NoArguments *Argument
)

func init() {
	NoArguments = &Argument{}
}

// Argument is an argument.
type Argument struct {
	// arg is the name of the argument.
	arg string
}

// ExactArgs returns an argument with the given name.
//
// Parameters:
//   - arg: The name of the argument.
//
// Returns:
//   - *Argument: The argument.
func ExactArgs(arg string) *Argument {
	return &Argument{
		arg: arg,
	}
}
