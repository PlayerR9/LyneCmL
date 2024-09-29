package cml

type RunFn func() error

type Command struct {
	Name  string
	RunFn RunFn
}
