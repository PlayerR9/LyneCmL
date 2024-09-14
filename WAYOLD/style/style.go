package style

import "github.com/gdamore/tcell"

// Enumer is an interface that can be enumerated.
// Value of 0 must be the background/text color.
type Enumer interface {
	~int

	// String returns the string representation of the enum.
	//
	// Returns:
	//   - string: The string representation of the enum.
	String() string
}

// Style is a style table.
type Style[T Enumer] struct {
	NormalText tcell.Style

	SuccessText tcell.Style

	ErrorText tcell.Style

	// table is the table of styles.
	table map[T]tcell.Style
}

// NewStyle returns a new style table.
//
// Returns:
//   - Style[T]: The new style table.
func NewStyle[T Enumer]() Style[T] {
	return Style[T]{
		NormalText:  tcell.StyleDefault,
		SuccessText: tcell.StyleDefault,
		ErrorText:   tcell.StyleDefault,
		table:       make(map[T]tcell.Style),
	}
}

// SetNormalText sets the normal text style.
//
// Parameters:
//   - style: The style to set.
func (s *Style[T]) SetNormalText(style tcell.Style) {
	if s == nil {
		return
	}

	s.NormalText = style
}

func (s *Style[T]) SetSuccessText(style tcell.Style) {
	if s == nil {
		return
	}

	s.SuccessText = style
}

func (s *Style[T]) SetErrorText(style tcell.Style) {
	if s == nil {
		return
	}

	s.ErrorText = style
}

// AddStyle adds a new style to the table.
//
// Parameters:
//   - name: The name of the style.
//   - style: The style to add.
//
// Does nothing if the receiver is nil. Moreover, if the style is already
// in the table, the new style will overwrite the old one.
func (s *Style[T]) AddStyle(name T, style tcell.Style) {
	if s == nil {
		return
	}

	s.table[name] = style
}
