package repl

import (
	"fmt"
	"labs/cli"
	"os"
	"reflect"
	"strings"
	"sync"
)

// Write buffer, an instruction buffer that writes its own contents to
// screen and current line
type wbuf []byte

// Process for wbuf. Actual implementation for writing the content
// in the wbuf
func (w wbuf) process(i *InOut) {
	if string(w) == "" {
		return
	}

	oldY := i.term.Cursor.Y

	i.write(w)

	parts := strings.Split(string(w), "\n")
	newLines := len(parts) - 1

	for _, v := range parts {
		i.term.Cursor.AddX(len(v))
	}

	i.term.Cursor.AddY(newLines)
	newY := i.term.Cursor.Y

	if oldY != newY {
		i.AddLines(newLines)
	}
}

// Read buffer, a buffer to hold Characters read in that need to
// be written
type rbuf []byte

// Process for rbuf. Actual implementation for rbuf
func (r rbuf) process(i *InOut) {
	if string(r) == "" {
		return
	}

	i.Wbuf = append(i.Wbuf, r...)
}

// Special buffer, an instruction buffer to hold special keystrokes
type spbuf []byte

// Process for spbuf. Actual implementation for the instructions
// in the spbuf
func (sp spbuf) process(i *InOut) {
	switch string(sp) {
	case "HOME":
		i.term.Cursor.Home(len(LINELOGO))
	case "END":
		i.term.Cursor.End(len(i.lines[i.term.Cursor.Y]))
	case "BACK":
		if i.term.Cursor.X <= len(LINELOGO) {
			return
		}

		i.lines[i.term.Cursor.Y] = i.lines[i.term.Cursor.Y].Backspace(i.term.Cursor.X)
		i.term.Cursor.AddX(-1)
	case "DEL":
		if i.term.Cursor.X <= len(LINELOGO) {
			return
		}
		i.lines[i.term.Cursor.Y] = i.lines[i.term.Cursor.Y].DelChar(i.term.Cursor.X)
	case "NEWL":
		i.term.Cursor.MoveTo(0, len(i.lines))
		i.term.Cursor.AddY(len(i.lines) - i.term.Cursor.Y)
		printLineLogo(i)

		i.term.Cursor.MoveTo(len(LINELOGO), len(i.lines))
		i.term.Cursor.X = len(LINELOGO)
		i.AddLines(1)

		i.done <- event{}
		return
	case "TAB":
		if i.term.Cursor.X >= i.term.Cols-8 || i.term.Cursor.X < 0 {
			return
		}

		i.lines[i.term.Cursor.Y] = i.lines[i.term.Cursor.Y].Tab(i.term.Cursor.X)
		i.term.Cursor.AddX(8)
	}
}

// Movement buffer, an instruction buffer to hold arrow keystrokes
type mvbuf []byte

// Process for mvbuf. Actual implementation for the instructions
// in the mvbuf
func (mv mvbuf) process(i *InOut) {
	switch string(mv) {
	case "UP":
		if i.term.Cursor.Y == 0 {
			return
		}

		i.term.Cursor.Up()
		i.term.Cursor.AddY(-1)
	case "DOWN":
		if i.term.Cursor.Y >= len(i.lines)-1 {
			return
		}

		i.term.Cursor.Down()
		i.term.Cursor.AddY(1)
	case "RIGHT":
		if i.term.Cursor.X == len(i.lines[i.term.Cursor.Y]) {
			return
		}

		i.term.Cursor.Right()
		i.term.Cursor.AddX(1)
	case "LEFT":
		if i.term.Cursor.X == len(LINELOGO) {
			return
		}

		i.term.Cursor.Left()
		i.term.Cursor.AddX(-1)
	}
}

// Filter buffer, input is placed here to be filtered and processed
// into instructions for the input buffer system.
type fbuf []byte

// Process for fbuf.
func (f fbuf) process(i *InOut) {
	f.filterInput(i)
}

// actuall concurrent filtering of fbuf if necessary.
func (f fbuf) filterInput(i *InOut) {
	var w sync.WaitGroup
	done := make(chan struct{}, 0)

	wg := &w
	wg.Add(5)

	go func() {
		go killCheck(i, wg)
		go testCheck(i, wg)
		go parseArrows(i, wg)
		go otherSpecial(i, wg)
		go regularChars(i, wg)

		wg.Wait()
		done <- event{}
	}()

	select {
	case <-done:
		close(done)
	}
}

// Filter through in!put byte for "Special" KeyStrokes: NL, CR, Home,
// End, Del, etc.
func otherSpecial(i *InOut, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 {
		wg.Done()
		return
	}

	switch string(i.Fbuf[0]) {
	case "\033":
		switch string(i.Fbuf[1]) {
		case "[":
			switch string(i.Fbuf[2]) {
			case "3":
				i.Spbuf = spbuf("DEL")
			case "H":
				i.Spbuf = spbuf("HOME")
			case "F":
				i.Spbuf = spbuf("END")
			}
		}
	case "\x7f":
		i.Spbuf = spbuf("BACK")
	case "\x0a", "\x0d":
		i.Spbuf = spbuf("NEWL")
	case "\x09":
		i.Spbuf = spbuf("TAB")
	default:
		wg.Done()
		return
	}

	i.Rbuf = rbuf("")

	wg.Done()
}

// Filter through in!put byte for "Movement" KeyStrokes: Arrows.
func parseArrows(i *InOut, wg *sync.WaitGroup) {
	if len(i.Fbuf) < 3 {
		wg.Done()
		return
	}

	switch string(i.Fbuf[1]) {
	case "[":
		switch string(i.Fbuf[2]) {
		case "A":
			i.Mvbuf = mvbuf("UP")
		case "B":
			i.Mvbuf = mvbuf("DOWN")
		case "C":
			i.Mvbuf = mvbuf("RIGHT")
		case "D":
			i.Mvbuf = mvbuf("LEFT")
		}
	default:
		wg.Done()
		return
	}

	i.Rbuf = rbuf("")

	wg.Done()
}

// Filter through in!put byte for "Quit/Kill" KeyStroke: Ctrl-C.
func killCheck(i *InOut, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x03" {
		wg.Done()
		return
	}

	i.term.Normal()
	cli.Restore()
	os.Exit(3)
}

// Filter through in!put byte for "Debug" KeyStroke: Ctrl-D.
func testCheck(i *InOut, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x04" {
		wg.Done()
		return
	}

	fmt.Printf("\n\rBUFFERS -- rbuf: %s|wbuf: %s|mvbuf: %s|spbuf: %s|fbuf: %s\n\r",
		string(i.Rbuf),
		string(i.Wbuf),
		string(i.Mvbuf),
		string(i.Spbuf),
		string(i.Fbuf),
	)

	fmt.Printf("Cursor -- X: %d|Y: %d\n\r",
		i.term.Cursor.X,
		i.term.Cursor.Y,
	)

	fmt.Printf("Lines -- count: %d| contents:\n\r",
		len(i.lines),
	)

	for _, v := range i.lines {
		fmt.Println(string(v))
	}

	wg.Done()
}

// Filter through in!put byte for "Regular" KeyStrokes: Characters,
// Spacebar.
func regularChars(i *InOut, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 {
		wg.Done()
		return
	}

	switch string(i.Fbuf[0]) {
	case "\x1b", "\x0a", "\x0d", "\x03", "\x04", "\x7f", "\x09":
		wg.Done()
		i.Rbuf = rbuf("")
		return
	default:
		i.Rbuf = append(i.Rbuf, i.Fbuf...)
		wg.Done()
	}
}

// Used for referring to all buffers as one entity.
type buffer interface {
	process(*InOut)
}

// Process a slice of buffers and resetting their value back to
// empty, individually.
func ProccessBuffers(bufs []buffer, in *InOut, wg *sync.WaitGroup) {

	for _, v := range bufs {
		v.process(in)

		reflect.ValueOf(v).Elem().SetBytes(
			[]byte(""),
		)
	}

	wg.Done()
	return
}
