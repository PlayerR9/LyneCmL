package screen

import (
	"context"
	"errors"
	"fmt"

	"github.com/gdamore/tcell"
)

type contextKey struct{}

type ContextOption func(c *Context)

func WithDefaultStyle(style tcell.Style) ContextOption {
	return func(c *Context) {
		c.screen.bg_style = style
	}
}

type Context struct {
	default_style tcell.Style

	screen *Screen
}

func NewContext(opts ...ContextOption) context.Context {
	c := &Context{
		default_style: tcell.StyleDefault,
	}

	for _, opt := range opts {
		opt(c)
	}

	ctx := context.WithValue(context.Background(), contextKey{}, c)

	return ctx
}

func from_context(ctx context.Context) *Context {
	return ctx.Value(contextKey{}).(*Context)
}

func Run(ctx context.Context) error {
	c := from_context(ctx)

	if c == nil {
		return errors.New("invalid context")
	}

	if c.screen != nil {
		return errors.New("screen already running")
	}

	screen, err := NewScreen(c.default_style)
	if err != nil {
		return fmt.Errorf("failed to create screen: %w", err)
	}

	c.screen = screen

	err = c.screen.Start()
	if err != nil {
		return fmt.Errorf("failed to start screen: %w", err)
	}

	return nil
}
