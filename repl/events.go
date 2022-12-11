package repl

type event struct{}

type Event interface {
	event | bool | struct{}
}

func EventChan(n ...int) chan struct{} {
	if len(n) > 1 {
		panic("cant add mult buffers to chan")
	}
	if len(n) != 0 {
		return make(chan struct{}, n[0])
	}

	return make(chan struct{})
}
