package display

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	// fs "github.com/PlayerR9/MyGoLib/Formatting/Strings"
	sfb "github.com/PlayerR9/safe/buffer"
)

///////////////////////////////////////////////////////

// Display is a display that can be used to display messages.
type Display struct {
	// buffer is the message buffer.
	buffer *sfb.Buffer[Msger]

	// history is a list of all the messages and interactions that have occurred
	// in the display since it started or since the last clear.
	history *sfb.SafeQueue[string]

	// wg is the wait group for the display.
	wg sync.WaitGroup

	// ctx is the context of the display.
	ctx context.Context

	// cancel is the cancel function of the display.
	cancel context.CancelFunc

	// configs is the configurations of the display.
	configs *DisplayConfigs

	// logger is the logger of the display.
	logger *log.Logger
}

// NewDisplay creates a new display with the given configurations.
//
// Parameters:
//   - config: The configurations of the display.
//
// Returns:
//   - *Display: The new display.
func NewDisplay(config *DisplayConfigs, logger *log.Logger) *Display {
	ctx, cancel := context.WithCancel(context.Background())

	buffer := sfb.NewBuffer[Msger]()
	history := sfb.NewSafeQueue[string]()

	return &Display{
		ctx:     ctx,
		cancel:  cancel,
		configs: config,
		buffer:  buffer,
		history: history,
		logger:  logger,
	}
}

// Start starts the display.
func (d *Display) Start() {
	d.buffer.Start()

	d.wg.Add(1)

	go d.msgListener()
}

// Close closes the display.
func (d *Display) Close() {
	d.buffer.Close()

	d.cancel()
	d.wg.Wait()
}

// msgListener listens for messages from the buffer.
func (d *Display) msgListener() {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			msg, ok := d.buffer.Receive()
			if !ok {
				return
			}

			err := d.msgHandler(msg)
			if err != nil {
				fmt.Println(err.Error())
				d.cancel()
				return
			}
		}
	}
}

// msgHandler handles messages from the buffer.
//
// Parameters:
//   - msg: The message to handle.
//
// Returns:
//   - error: An error if the message failed to be handled.
func (d *Display) msgHandler(msg any) error {
	switch msg := msg.(type) {
	case *TextMsg:
		// space := strings.Repeat(" ", d.configs.Spacing)

		// str, _ := fs.FixTabStop(0, d.configs.TabSize, space, msg.text)

		// d.history.Enqueue(str)

		// fmt.Println(str)
	case *LogMsg:
		// space := strings.Repeat(" ", d.configs.Spacing)

		// str, _ := fs.FixTabStop(0, d.configs.TabSize, space, msg.text)

		// d.history.Enqueue(str)

		// d.logger.Println(str)
	case *ClearHistoryMsg:
		d.history.Clear()
	case *StoreHistoryMsg:
		file, err := os.Create(msg.loc)
		if err != nil {
			return err
		}
		defer file.Close()

		for _, line := range d.history.Slice() {
			_, err := file.WriteString(line + "\n")
			if err != nil {
				return err // TODO: Use ErrWhile
			}
		}
	case *AbruptExitMsg:
		fmt.Println("History:")
		fmt.Println()

		for _, line := range d.history.Slice() {
			fmt.Println(line)
		}

		return msg.reason
	case *InputMsg:
		if msg.text != "" {
			fmt.Println(msg.text)
		}

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

// Send sends a message to the display.
//
// Parameters:
//   - msg: The message to send.
//
// Returns:
//   - error: An error if the message failed to send.
//
// Behaviors:
//   - If the message is nil, nothing will be sent.
//   - It can only error if the context is done.
func (d *Display) Send(msg any) error {
	if msg == nil {
		return nil
	}

	select {
	case <-d.ctx.Done():
		return d.ctx.Err()
	default:
		d.buffer.Send(msg)
		return nil
	}
}

// IsDone returns whether the display is done.
//
// Returns:
//   - bool: Whether the display is done.
func (d *Display) IsDone() bool {
	select {
	case <-d.ctx.Done():
		return true
	default:
		return false
	}
}
