package assert

// ErrAssertFail is an error that is returned when an assertion fails.
type ErrAssertFail struct {
	// Msg is the message of the error.
	Msg string
}

// Error implements the error interface.
//
// Message:
//
//	[ASSERT FAILED]: <msg>
//
// If 'msg' is empty, 'an assertion has failed' is used instead.
func (err ErrAssertFail) Error() string {
	msg := "[ASSERT FAILED]: "

	if err.Msg == "" {
		msg += "an assertion has failed"
	} else {
		msg += err.Msg
	}

	return msg
}

// NewErrAssertFail returns a new ErrAssertFail.
//
// Parameters:
//   - msg: the message of the error.
//
// Returns:
//   - *ErrAssertFail: the new error. Never returns nil.
func NewErrAssertFail(msg string) *ErrAssertFail {
	return &ErrAssertFail{
		Msg: msg,
	}
}
