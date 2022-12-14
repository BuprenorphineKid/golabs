package repl

/* I realize alot of the implementation and logic in this file can
be a little bit hard to follow.. Its literally all because i
couldnt let go of this little clever bit of dynamically identifying
the line youre on by using the i.term.Cursor's Y position as the
array index. Like Terminal Cursor Positioning is already 1 indexed
if you get it by means of escape sequences like i did. Might as
well put it to a little use lol. Anyway, just keep that in mind
and it really isnt all that bad*/

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
		i.term.Cursor.X = len(LINELOGO)
	case "END":
		i.term.Cursor.End(len(i.lines[i.term.Cursor.Y]) + len(LINELOGO))
		i.term.Cursor.X = len(i.lines[i.term.Cursor.Y]) + len(LINELOGO)
	case "BACK":
		if i.term.Cursor.X <= len(LINELOGO) ||
			i.term.Cursor.X-len(LINELOGO) > len(i.lines[i.term.Cursor.Y]) {
			return
		}

		i.lines[i.term.Cursor.Y] = i.lines[i.term.Cursor.Y].Backspace(i.term.Cursor.X)
		i.term.Cursor.AddX(-1)

		i.term.Cursor.Left()

		Refresh(i)

		i.term.Cursor.Left()
	case "DEL":
		if i.term.Cursor.X < len(LINELOGO) {
			return
		}

		i.lines[i.term.Cursor.Y] = i.lines[i.term.Cursor.Y].DelChar(i.term.Cursor.X)

		Refresh(i)
		i.term.Cursor.Left()
	case "NEWL":
		i.term.Cursor.MoveTo(0, len(i.lines))
		i.term.Cursor.AddY(len(i.lines) - i.term.Cursor.Y)

		i.AddLines(1)

		i.done <- event{}
		return
	case "TAB":
		if i.term.Cursor.X >= i.term.Cols-8 ||
			i.term.Cursor.X < 0 {
			return
		}

		i.lines[i.term.Cursor.Y] = i.lines[i.term.Cursor.Y].Tab(i.term.Cursor.X)
		i.term.Cursor.AddX(8)

		Refresh(i)
	}
}

// Movement buffer, an instruction buffer to hold arrow keystrokes
type mvbuf []byte

// Process for mvbuf. Actual implementation for the instructions
// in the mvbuf
func (mv mvbuf) process(i *InOut) {
	switch string(mv) {
	case "UP":
		if i.term.Cursor.Y <= 5 {
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
		if i.term.Cursor.X >= len(LINELOGO)+len(i.lines[i.term.Cursor.Y]) {
			return
		}

		i.term.Cursor.Right()
		i.term.Cursor.AddX(1)
	case "LEFT":
		if i.term.Cursor.X <= len(LINELOGO) {
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

// actual concurrent filtering of fbuf as necessary.
func (f fbuf) filterInput(i *InOut) {
	var w sync.WaitGroup
	done := make(chan struct{}, 0)

	wg := &w
	wg.Add(5)

	go func() {
		go killCheck(i, wg)
		go DebugCheck(i, wg)
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

// Filter through in!put bytes for "Special" KeyStrokes: NL, CR, Home,
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

// Filter through in!put bytes for "Movement" KeyStrokes: Arrows.
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

// Filter through input byte for "Debug" KeyStroke: Ctrl-D.
func DebugCheck(i *InOut, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x04" {
		wg.Done()
		return
	}

	i.term.Cursor.SavePos()
	old := i.term.Cursor.Y

	x := i.term.Cols / 3

	i.term.Cursor.MoveTo(x+x, 1)

	fmt.Printf("POS - |X: %d| |Y: %d|",
		i.term.Cursor.X,
		i.term.Cursor.Y,
	)

	i.term.Cursor.MoveTo(x+x, 2)

	fmt.Printf("LINES - |count: %d|",
		len(i.lines),
	)

	i.term.Cursor.MoveTo(x+x, 3)
	i.term.Cursor.Y = 3
	for n, v := range i.lines[:] {
		fmt.Printf("%d - |%s|\n", n, string(v))
		i.term.Cursor.AddY(1)

		i.term.Cursor.MoveTo(x+x, i.term.Cursor.Y)
	}

	i.term.Cursor.RestorePos()
	i.term.Cursor.Y = old

	wg.Done()
}

// Filter through input byte for "Regular" KeyStrokes: Characters,
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
