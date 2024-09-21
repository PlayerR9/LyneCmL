package simple

import (
	"fmt"
	"os"
)

// DefaultExitSequence is a method that prints the given error and exits the program.
//
// Parameters:
//   - err: The error to print.
func DefaultExitSequence(err error) {
	var exit_code int

	if err == nil {
		_, err := fmt.Println("Command ran successfully")
		if err != nil {
			panic(err)
		}

		exit_code = 0
	} else {
		_, err := fmt.Println(err.Error())
		if err != nil {
			panic(err)
		}

		exit_code = 1
	}

	_, err = fmt.Println()
	if err != nil {
		panic(err)
	}

	_, err = fmt.Println("Press ENTER to exit...")
	if err != nil {
		panic(err)
	}

	_, err = fmt.Scanln()
	if err != nil {
		panic(err)
	}

	os.Exit(exit_code)
}
