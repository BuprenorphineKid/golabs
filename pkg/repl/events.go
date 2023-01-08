package repl

/*
making BackGround Events, namely, evaluator, who is constantly running the
beneath go script for Instant results ready, on the fly. rather than having to
type the ";eval command. the event is passed in two chans a done chan, and
results chan, results is what gets returned upon calling the Ping method.
neatly structured into an anonymous struct"
*/

type event struct{}

func EventChan(n ...int) chan struct{} {
	if len(n) > 1 {
		panic("cant add mult buffers to chan")
	}
	if len(n) != 0 {
		return make(chan struct{}, n[0])
	}

	return make(chan struct{})
}
