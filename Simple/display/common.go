package display

type Displayer interface {
	Send(msg any)
}
