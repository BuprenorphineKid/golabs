package readline

// Custom string type to represent the inner line buffer that gets
// sent around and evaluated and what not. The lines should always
// directly mimic the lines on the screen.
type line string

// Line method for deleting a character behind the cursor when the
// the Backspace button is preased.
func (l *line) Backspace(xpos int) line {
	if xpos-1 <= 0 {
		return *l
	}

	ch := make(chan line)

	go func() {
		pos := (xpos - 1) - len(INPROMPT)
		b := []byte(*l)

		front := b[:pos]
		back := b[pos+1:]

		ch <- line(append(front, back...))
		close(ch)
	}()

	return <-ch
}

// Line method used to delete a character that the cursor is highliting
// whenever the del button is pressed.
func (l *line) DelChar(xpos int) line {
	if (xpos-1)-len(INPROMPT)+1 > len(*l)-1 {
		return *l
	}

	if (xpos-1)-len(INPROMPT)+1 == len(*l)-1 {
		if len(*l) == 1 {
			return line("")
		}
		var b = []byte(*l)

		return line(b[:len(b)-2])
	}
	var pos int

	if (xpos-1)-len(INPROMPT) < 0 {
		pos = 0
	} else {
		pos = (xpos - 1) - len(INPROMPT)
	}

	ch := make(chan line)

	go func() {
		b := []byte(*l)

		front := b[:pos+1]
		var back = b[pos+2:]

		ch <- line(append(front, back...))
		close(ch)
	}()

	return <-ch
}

// Line method to margin the text. done with 8 spaces rather than a
// single tab, strictly due to convenience.
func (l *line) Tab(xpos int) line {
	const tab = "    "
	var pos int

	if xpos-1-len(INPROMPT) <= 0 {
		pos = 0
	} else {
		pos = xpos - 1 - len(INPROMPT)
	}

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

// Use this method if you want to insert a character into the text
// inside the line without overwriting other pre-existing
// characters.
func (l *line) Insert(char []byte, xpos int) line {
	var pos int

	if (xpos-1)-len(INPROMPT) <= 0 {
		return l.Write(char, xpos)
	} else {
		pos = (xpos - 1) - len(INPROMPT)
	}

	ch := make(chan line)

	go func() {
		b := []byte(*l)

		front := b[:pos+1]
		back := b[pos+1:]
		mod := append(char, back...)

		ch <- line(append(front, mod...))
		close(ch)
	}()

	return <-ch
}

// Only use this method for the start of a line where deciding the
// cutting positions would cause a panic.
func (l *line) Write(char []byte, xpos int) line {
	if (xpos-1)-len(INPROMPT) > 0 {
		return l.Insert(char, xpos)
	}

	return line(append([]byte(*l), char...))
}
