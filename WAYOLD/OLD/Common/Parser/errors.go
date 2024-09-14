package Parser

import (
	"strings"

	cms "github.com/PlayerR9/LyneCml/OLD/Simple"
)

type ErrFlagNotFound struct {
	IsShort bool
	Name    string
}

func (e *ErrFlagNotFound) Error() string {
	var builder strings.Builder

	builder.WriteString("flag ")

	if e.IsShort {
		builder.WriteRune('-')
	} else {
		builder.WriteString("--")
	}

	builder.WriteString(e.Name)

	builder.WriteString(" not found")

	return builder.String()
}

func NewErrFlagNotFound(isShort bool, name string) *ErrFlagNotFound {
	e := &ErrFlagNotFound{
		IsShort: isShort,
		Name:    name,
	}
	return e
}

type ErrFlagMissingArg struct {
	IsShort bool
	Flag    *cms.Flag
}

func (e *ErrFlagMissingArg) Error() string {
	var builder strings.Builder

	builder.WriteString("flag ")

	if e.IsShort {
		builder.WriteRune('-')
		builder.WriteRune(e.Flag.ShortName)
	} else {
		builder.WriteString("--")
		builder.WriteString(e.Flag.LongName)
	}

	builder.WriteString(" requires an argument")

	return builder.String()
}

func NewErrFlagMissingArg(isShort bool, flag *cms.Flag) *ErrFlagMissingArg {
	e := &ErrFlagMissingArg{
		IsShort: isShort,
		Flag:    flag,
	}
	return e
}

type ErrInvalidFlag struct {
	IsShort bool
	Flag    *cms.Flag
	Reason  error
}

func (e *ErrInvalidFlag) Error() string {
	var builder strings.Builder

	builder.WriteString("flag ")

	if e.IsShort {
		builder.WriteRune('-')
		builder.WriteRune(e.Flag.ShortName)
	} else {
		builder.WriteString("--")
		builder.WriteString(e.Flag.LongName)
	}

	builder.WriteString(" is invalid")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

func (e *ErrInvalidFlag) Unwrap() error {
	return e.Reason
}

func (e *ErrInvalidFlag) ChangeReason(reason error) {
	e.Reason = reason
}

func NewErrInvalidFlag(isShort bool, flag *cms.Flag, reason error) *ErrInvalidFlag {
	e := &ErrInvalidFlag{
		IsShort: isShort,
		Flag:    flag,
		Reason:  reason,
	}
	return e
}

type ErrInvalidArg struct {
	Arg     string
	IsShort bool
	Reason  error
}

func (e *ErrInvalidArg) Error() string {
	var builder strings.Builder

	if e.IsShort {
		builder.WriteRune('-')
	} else {
		builder.WriteString("--")
	}

	builder.WriteString(e.Arg)
	builder.WriteString(" is invalid")

	if e.Reason != nil {
		builder.WriteString(": ")
		builder.WriteString(e.Reason.Error())
	}

	return builder.String()
}

func (e *ErrInvalidArg) Unwrap() error {
	return e.Reason
}

func (e *ErrInvalidArg) ChangeReason(reason error) {
	e.Reason = reason
}

func NewErrInvalidArg(arg string, isShort bool, reason error) *ErrInvalidArg {
	e := &ErrInvalidArg{
		Arg:     arg,
		IsShort: isShort,
		Reason:  reason,
	}
	return e
}
