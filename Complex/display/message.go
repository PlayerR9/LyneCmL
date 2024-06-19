package display

import (
	"errors"
	"fmt"
)

type Msger interface{}

// TextMsg is a message that contains text.
type TextMsg struct {
	// text is the text of the message.
	text string
}

// NewTextMsg creates a new TextMsg.
//
// It must not end with a newline character.
//
// Parameters:
//   - text: The text of the message.
//
// Returns:
//   - *TextMsg: The new TextMsg.
func NewTextMsg(text string) *TextMsg {
	return &TextMsg{
		text: text,
	}
}

// ClearHistoryMsg is a message that clears the history of the display.
type ClearHistoryMsg struct{}

// NewClearHistoryMsg creates a new ClearHistoryMsg.
//
// Returns:
//   - *ClearHistoryMsg: The new ClearHistoryMsg.
func NewClearHistoryMsg() *ClearHistoryMsg {
	return &ClearHistoryMsg{}
}

// StoreHistoryMsg is a message that makes a backup of the history.
type StoreHistoryMsg struct {
	// loc is the directory to store the history in.
	loc string
}

// NewStoreHistoryMsg creates a new StoreHistoryMsg.
//
// Parameters:
//   - loc: The directory to store the history in.
//
// Returns:
//   - *StoreHistoryMsg: The new StoreHistoryMsg.
//
// Behaviors:
//   - If loc is an empty string, nil will be returned.
func NewStoreHistoryMsg(loc string) *StoreHistoryMsg {
	if loc == "" {
		return nil
	}

	return &StoreHistoryMsg{
		loc: loc,
	}
}

// AbruptExitMsg is a message that causes the display to abruptly exit.
type AbruptExitMsg struct {
	// reason is the error that caused the abrupt exit.
	reason error
}

// NewAbruptExitMsg creates a new AbruptExitMsg.
//
// Parameters:
//   - reason: The error that caused the abrupt exit.
//
// Returns:
//   - *AbruptExitMsg: The new AbruptExitMsg.
//
// Behaviors:
//   - If reason is nil, it will be set to the error "no reason provided".
func NewAbruptExitMsg(reason error) *AbruptExitMsg {
	if reason == nil {
		reason = errors.New("no reason provided")
	}

	return &AbruptExitMsg{
		reason: reason,
	}
}

// InputMsg is a message that requests input from the user.
type InputMsg struct {
	// text is the text to display to the user.
	text string

	// receiveCh is the channel to receive the input on.
	receiveCh chan any
}

// NewInputMsg creates a new InputMsg.
//
// Parameters:
//   - text: The text to display to the user.
//
// Returns:
//   - *InputMsg: The new InputMsg.
func NewInputMsg(text string) *InputMsg {
	ch := make(chan any)

	return &InputMsg{
		text:      text,
		receiveCh: ch,
	}
}

// Receive receives the input from the user.
//
// Returns:
//   - string: The input from the user.
//   - error: An error if the input failed.
func (im *InputMsg) Receive() (string, error) {
	defer close(im.receiveCh)

	input, ok := <-im.receiveCh
	if !ok {
		return "", errors.New("input channel closed")
	}

	switch input := input.(type) {
	case error:
		return "", input
	case string:
		return input, nil
	default:
		return "", fmt.Errorf("unexpected input type: %T", input)
	}
}
