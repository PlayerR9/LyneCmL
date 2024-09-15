package simple

import (
	"context"
	"errors"

	ds "github.com/PlayerR9/LyneCml/WAYOLD/screen"
	gcint "github.com/PlayerR9/go-commons/ints"
)

type contextKey struct{}

type Context struct {
	screen *ds.Screen
}

func NewContext() context.Context {
	c := &Context{}

	ctx := context.WithValue(context.Background(), contextKey{}, c)

	return ctx
}

func from_context(ctx context.Context) *Context {
	return ctx.Value(contextKey{}).(*Context)
}

func Do(ctx context.Context, actions ...Action) error {
	c := from_context(ctx)

	if c == nil {
		return errors.New("invalid context")
	}

	for i, action := range actions {
		if action == nil {
			continue
		}

		err := do_action(c, action)
		if err != nil {
			return gcint.NewErrWhileAt("doing", i+1, "action", err)
		}
	}

	return nil
}
