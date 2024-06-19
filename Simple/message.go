package Simple

import (
	"fmt"
	"path"
	"strings"

	fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	ufm "github.com/PlayerR9/MyGoLib/Utility/FileManager"
)

// msgHandler handles messages from the buffer.
//
// Parameters:
//   - msg: The message to handle.
//
// Returns:
//   - error: An error if the message failed to be handled.
func (p *Program) msgHandler(msg any) error {
	switch msg := msg.(type) {
	case *TextMsg:
		space := strings.Repeat(" ", p.Options.Spacing)
		str := msg.text

		str, _ = fs.FixTabStop(0, p.Options.TabSize, space, str)

		p.history.Enqueue(str)
		fmt.Print(str)
	case *ClearHistoryMsg:
		p.history.Clear()
	case *StoreHistoryMsg:
		iter := p.history.Iterator()

		fullpath := path.Join("partials", msg.filename)

		fw := ufm.NewFileWriter(fullpath)
		err := fw.Create()
		if err != nil {
			return err
		}
		defer fw.Close()

		for {
			line, err := iter.Consume()
			if err != nil {
				break
			}

			err = fw.AppendLine(line)
			if err != nil {
				return err
			}
		}
	case *AbruptExitMsg:
		fmt.Println("History:")
		fmt.Println()

		iter := p.history.Iterator()

		for {
			line, err := iter.Consume()
			if err != nil {
				break
			}

			fmt.Println(line)
		}

		return msg.err
	case *InputMsg:
		fmt.Println(msg.text)
		fmt.Print("> ")

		var input string

		_, err := fmt.Scanln(&input)
		if err != nil {
			msg.receiveCh <- err
		} else {
			msg.receiveCh <- input
		}
	default:
		return fmt.Errorf("unknown message type: %T", msg)
	}

	return nil
}

// TextMsg is a message that contains text.
type TextMsg struct {
	// text is the text of the message.
	text string
}

// Println prints a line to the program's buffer.
//
// Parameters:
//   - args: The items to print.
func (p *Program) Println(args ...interface{}) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
		str := fmt.Sprintln(args...)

		p.buffer.Send(&TextMsg{text: str})

		return nil
	}
}

// Printf prints a formatted line to the program's buffer.
//
// Parameters:
//   - format: The format of the line.
//   - args: The items to print.
//
// Behaviors:
//   - A newline character will be appended to the end of the string
//     if it does not already have one.
func (p *Program) Printf(format string, args ...interface{}) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
		str := fmt.Sprintf(format, args...)

		p.buffer.Send(&TextMsg{text: str + "\n"})

		return nil
	}
}

// ClearHistoryMsg is a message that clears the history of the program.
type ClearHistoryMsg struct{}

// ClearHistory clears the history of the program.
//
// This function is thread-safe.
func (p *Program) ClearHistory() error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
		p.buffer.Send(&ClearHistoryMsg{})
		return nil
	}
}

// StoreHistoryMsg is a message that stores the history of the program to a file.
type StoreHistoryMsg struct {
	// filename is the name of the file to store the history in.
	filename string
}

// SavePartial saves the current history to a file in the partials directory.
//
// This can be used for logging/debugging purposes and/or to save the state of
// the program or evaluate the program's output.
//
// This function is thread-safe.
//
// Parameters:
//   - filename: The name of the file to save the partial to.
func (p *Program) SavePartial(filename string) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
		p.buffer.Send(&StoreHistoryMsg{filename: filename})

		return nil
	}
}

// AbruptExitMsg is a message that causes the program to abruptly exit.
type AbruptExitMsg struct {
	// err is the error that caused the abrupt exit.
	err error
}

// Panic causes the program to abruptly exit with the given error.
//
// This function is thread-safe.
//
// Parameters:
//   - err: The error that caused the abrupt exit.
func (p *Program) Panic(err error) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
		if err == nil {
			err = fmt.Errorf("no reason given")
		}

		p.buffer.Send(&AbruptExitMsg{err: err})

		return nil
	}
}

// InputMsg is a message that requests input from the user.
type InputMsg struct {
	// text is the text to display to the user.
	text string

	// receiveCh is the channel to receive the input on.
	receiveCh chan<- any
}

// Input requests input from the user.
//
// This function is thread-safe.
//
// Parameters:
//   - text: The text to display to the user.
//
// Returns:
//   - string: The input from the user.
//   - error: An error if the input failed.
func (p *Program) Input(text string) (string, error) {
	select {
	case <-p.ctx.Done():
		return "", p.ctx.Err()
	default:
		receiveCh := make(chan any)
		defer close(receiveCh)

		p.buffer.Send(&InputMsg{
			text:      text,
			receiveCh: receiveCh,
		})

		input := <-receiveCh

		switch input := input.(type) {
		case string:
			return input, nil
		case error:
			return "", input
		default:
			return "", fmt.Errorf("unexpected input type: %T", input)
		}
	}
}
