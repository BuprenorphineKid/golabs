package repl

/* I realize alot of the implementation and logic in this file can
be a little bit hard to follow.. Its literally all because i
couldnt let go of this little clever bit of dynamically identifying
the line youre on by using the c's Y position as the
array index. Like Terminal Cursor Positioning is already 1 indexed
if you get it by means of escape sequences like i did. Might as
well put it to a little use lol. Anyway, just keep that in mind
and it really isnt all that bad*/

import (
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
func (w wbuf) process(i *Input, c Cursor) {
	if string(w) == "" {
		return
	}

	oldY := reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)

	i.write(w)

	parts := strings.Split(string(w), "\n")
	newLines := len(parts) - 1

	for _, v := range parts {
		c.AddX(len(v))
	}

	c.AddY(newLines)
	newY := reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)

	if oldY != newY {
		i.AddLines(newLines)
	}
}

// Read buffer, a buffer to hold Characters read in that need to
// be written
type rbuf []byte

// Process for rbuf. Actual implementation for rbuf
func (r rbuf) process(i *Input, c Cursor) {
	if string(r) == "" {
		return
	}

	i.Wbuf = append(i.Wbuf, r...)
}

// Special buffer, an instruction buffer to hold special keystrokes
type spbuf []byte

// Process for spbuf. Actual implementation for the instructions
// in the spbuf
func (sp spbuf) process(i *Input, c Cursor) {
	switch string(sp) {
	case "HOME":
		c.Home(len(INPROMPT))
		c.SetX(len(INPROMPT))
	case "END":
		c.End(len(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]) + len(INPROMPT))
		c.SetX(len(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]) + len(INPROMPT))
	case "BACK":
		if reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int) <= len(INPROMPT) ||
			reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int)-len(INPROMPT) > len(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]) {
			return
		}

		i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)] = i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)].Backspace(reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int))
		c.AddX(-1)

		c.Left()

		output.SetLine(string(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]))
		output.devices["main"].(Display).RenderLine(c)

		c.Left()
	case "DEL":
		if reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int) < len(INPROMPT) ||
			reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int)-len(INPROMPT) > len(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]) {
			return
		}

		i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)] = i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)].DelChar(reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int))

		output.SetLine(string(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]))
		output.devices["main"].(Display).RenderLine(c)
	case "NEWL":
		c.MoveTo(0, len(i.lines))
		c.AddY(len(i.lines) - reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int))
		c.SetX(0)

		i.AddLines(1)

		i.done <- event{}
		return
	case "TAB":
		if reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int) >= i.term.Cols-8 ||
			reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int) < 0 {
			return
		}

		i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)] = i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)].Tab(reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int))
		c.AddX(4)

		output.SetLine(string(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]))
		output.devices["main"].(Display).RenderLine(c)
	}
}

// Movement buffer, an instruction buffer to hold arrow keystrokes
type mvbuf []byte

// Process for mvbuf. Actual implementation for the instructions
// in the mvbuf
func (mv mvbuf) process(i *Input, c Cursor) {
	switch string(mv) {
	case "UP":
		if reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int) <= 5 {
			return
		}

		c.Up()
		c.AddY(-1)
	case "DOWN":
		if reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int) >= len(i.lines)-1 {
			return
		}

		c.Down()
		c.AddY(1)
	case "RIGHT":
		if reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int) >= len(INPROMPT)+len(i.lines[reflect.ValueOf(c).Elem().FieldByName("Y").Interface().(int)]) {
			return
		}

		c.Right()
		c.AddX(1)
	case "LEFT":
		if reflect.ValueOf(c).Elem().FieldByName("X").Interface().(int) <= len(INPROMPT) {
			return
		}

		c.Left()
		c.AddX(-1)
	}
}

// Filter buffer, input is placed here to be filtered and processed
// into instructions for the input buffer system.
type fbuf []byte

// Process for fbuf.
func (f fbuf) process(i *Input, c Cursor) {
	f.filterInput(i)
}

// actual concurrent filtering of fbuf as necessary.
func (f fbuf) filterInput(i *Input) {
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
func otherSpecial(i *Input, wg *sync.WaitGroup) {
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
func parseArrows(i *Input, wg *sync.WaitGroup) {
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

// Filter through input byte for "Quit/Kill" KeyStroke: Ctrl-C.
func killCheck(i *Input, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x03" {
		wg.Done()
		return
	}

	i.term.Normal()
	cli.Restore()
	os.Exit(3)
}

// Filter through input byte for "Debug" KeyStroke: Ctrl-B.
func DebugCheck(i *Input, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x02" {
		wg.Done()
		return
	}

	if i.InDebug {
		i.Debugger.Off <- event{}
		i.Debugger = new(Debugger)
		i.InDebug = false
	} else {
		i.Debugger = NewDebugger()

		go i.Debugger.DebugMode(i)

		i.InDebug = true
	}

	wg.Done()
}

// Filter through input byte for "Regular" KeyStrokes: Characters,
// Spacebar.
func regularChars(i *Input, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 {
		wg.Done()
		return
	}

	switch string(i.Fbuf[0]) {
	case "\x1b", "\x0a", "\x0d", "\x03", "\x04", "\x7f", "\x09", "\x11",
		"\x12", "\x01", "\x02", "\x05", "\x06", "\x07", "\x08", "\x13":
		//		"\x22", "\x23", "\x24", "\x25", "\x26":
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
	process(*Input, Cursor)
}

// Process a slice of buffers and resetting their value back to
// empty, individually.
func ProccessBuffers(bufs []buffer, i *Input, c Cursor, wg *sync.WaitGroup) {

	for _, v := range bufs {
		v.process(i, c)

		reflect.ValueOf(v).Elem().SetBytes(
			[]byte(""),
		)
	}

	wg.Done()
	return
}
