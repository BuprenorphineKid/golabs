package repl

type line string

func (l *line) Backspace() {
	print("\b \b")
}

func (l *line) DelChar() {
	print(" \b")
}
