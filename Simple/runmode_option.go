package Simple

import (
	"fmt"
	"os"
	"sync"

	sfb "github.com/PlayerR9/MyGoLib/Safe/Buffer"
	ufm "github.com/PlayerR9/MyGoLib/Utility/FileManager"
)

// RunModeOption is a function that modifies a Program.
type RunModeOption func(*programOptions)

// WithMsgLocation sets the location of print messages.
//
// Parameters:
//   - msgLocation: The name of the file where messages will be printed.
//
// Returns:
//   - RunModeOption: The option to set the location of print messages.
//
// Behaviors:
//   - If msgLocation is empty, the location will be set to "stdout".
func WithMsgLocation(msgLocation string) RunModeOption {
	return func(p *programOptions) {
		p.messageReceiverLocation = msgLocation
	}
}

type programOptions struct {
	messageReceiverListener func()
	messageReceiverLocation string

	wg sync.WaitGroup

	buffer *sfb.Buffer[string]
}

func newProgramOptions(buffer *sfb.Buffer[string], opts ...RunModeOption) *programOptions {
	p := &programOptions{
		messageReceiverListener: func() {},
		messageReceiverLocation: "",
		buffer:                  buffer,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *programOptions) startMsgListener() error {
	if p.messageReceiverLocation == "" {
		go func() {
			defer p.wg.Done()

			for {
				str, ok := p.buffer.Receive()
				if !ok {
					break
				}

				fmt.Print(str)
			}
		}()
	} else {
		fw := ufm.NewFileWriter(p.messageReceiverLocation, os.ModePerm, 0666)

		exists, err := fw.Exists()
		if err != nil {
			return fmt.Errorf("could not check if file exists: %w", err)
		}

		if !exists {
			err = fw.Create()
			if err != nil {
				return fmt.Errorf("could not create file: %w", err)
			}
		}

		go func() {
			defer p.wg.Done()

			for {
				str, ok := p.buffer.Receive()
				if !ok {
					break
				}

				err := fw.AppendLine(str)
				if err != nil {
					fmt.Println("could not write to file:", err)

					break
				}
			}
		}()
	}

	return nil
}
