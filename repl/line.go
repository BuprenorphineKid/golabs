package repl

type line string

func (l *line) Backspace(xpos int) line {
	if xpos-1 <= 0 {
		return *l
	}

	print("\b \b")

	ch := make(chan line)

	go func() {
		pos := (xpos - 1) - len(LINELOGO)
		b := []byte(*l)

		front := b[:pos]
		back := b[pos+1:]

		ch <- line(append(front, back...))
		close(ch)
	}()

	return <-ch
}

func (l *line) DelChar(xpos int) line {
	if xpos-1-len(LINELOGO) <= 0 || (xpos-1)-len(LINELOGO)+1 > len(*l)-1 {
		return *l
	}

	print(" \b")

	ch := make(chan line)

	go func() {
		pos := (xpos - 1) - len(LINELOGO)
		b := []byte(*l)

		front := b[:pos+1]
		var back = b[pos+1:]
		mod := append([]byte(" "), back...)

		ch <- line(append(front, mod...))
		close(ch)
	}()

	return <-ch
}

func (l *line) Tab(xpos int) line {
	const tab = "        "
	var pos int

	if (xpos-1)-len(LINELOGO) <= 0 {
		pos = 0
	} else {
		pos = (xpos - 1) - len(LINELOGO)
	}

	print(tab)

	ch := make(chan line)

	go func() {
		b := []byte(*l)

		front := b[:pos]
		back := b[pos:]
		mod := append([]byte(tab), back...)

		ch <- line(append(front, mod...))
		close(ch)
	}()

	return <-ch
}
