package parsing

import (
	"fmt"
	"regexp"
)

type Tokener interface {
	~int

	fmt.Stringer
}

type RegexTable[T Tokener] struct {
	regexes []*regexp.Regexp
	tokens  []T
}

func NewRegexTable[T Tokener]() *RegexTable[T] {
	return &RegexTable[T]{
		regexes: make([]*regexp.Regexp, 0),
		tokens:  make([]T, 0),
	}
}

func (t *RegexTable[T]) AddRegex(t_type T, expr string) error {
	compiled, err := regexp.Compile("^" + expr)
	if err != nil {
		return err
	}

	t.regexes = append(t.regexes, compiled)
	t.tokens = append(t.tokens, t_type)

	return nil
}

func (t *RegexTable[T]) MustAddRegex(t_type T, expr string) {
	compiled := regexp.MustCompile("^" + expr)
	t.regexes = append(t.regexes, compiled)
	t.tokens = append(t.tokens, t_type)
}
