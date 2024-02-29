package readline

/* I realize alot of the implementation and logic in this file can
be a little bit hard to follow.. Its literally all because i
couldnt let go of this little clever bit of dynamically identifying
the line youre on by using the c's Y position as the
array index. Like Terminal Cursor Positioning is already 1 indexed
if you get it by means of escape sequences like i did. Might as
well put it to a little use lol. Anyway, just keep that in mind
and it really isnt all that bad*/

import (
	"labs/pkg/cli"
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
func (w wbuf) process(i *Input) {
	if string(w) == "" {
		return
	}

	oldY := Term.Cursor.GetY()

	i.write(w)

	parts := strings.Split(string(w), "\n")
	newLines := len(parts) - 1

	for _, v := range parts {
		Term.Cursor.AddX(len(v))
	}

	Term.Cursor.AddY(newLines)
	newY := Term.Cursor.GetY()

	if oldY != newY {
		i.AddLines(newLines)
	}
}

// Read buffer, a buffer to hold Characters read in that need to
// be written
type rbuf []byte

// Process for rbuf. Actual implementation for rbuf
func (r rbuf) process(i *Input) {
	if string(r) == "" {
		return
	}

	i.Wbuf = append(i.Wbuf, r...)
}

// Special buffer, an instruction buffer to hold special keystrokes
type spbuf []byte

// Process for spbuf. Actual implementation for the instructions
// in the spbuf
func (sp spbuf) process(i *Input) {
	switch string(sp) {
	case "HOME":
		Term.Cursor.Home(len(INPROMPT))
		Term.Cursor.SetX(len(INPROMPT))
	case "END":
		Term.Cursor.End(len(i.Lines[len(i.Lines)-1]) + len(INPROMPT))
		Term.Cursor.SetX(len(i.Lines[len(i.Lines)-1]) + len(INPROMPT))
	case "BACK":
		if Term.Cursor.GetX() <= len(INPROMPT) ||
			Term.Cursor.GetX()-len(INPROMPT) > len(i.Lines[len(i.Lines)-1]) {
			return
		}

		i.Lines[len(i.Lines)-1] = i.Lines[len(i.Lines)-1].Backspace(Term.Cursor.GetX())
		Term.Cursor.AddX(-1)

		Term.Cursor.Left()

		out.SetLine(string(i.Lines[len(i.Lines)-1]))
		out.Devices["main"].(Display).RenderLine()

		Term.Cursor.Left()
	case "DEL":
		if Term.Cursor.GetX() < len(INPROMPT) ||
			Term.Cursor.GetX()-len(INPROMPT) > len(i.Lines[len(i.Lines)-1]) {
			return
		}

		i.Lines[len(i.Lines)-1] = i.Lines[len(i.Lines)-1].DelChar(Term.Cursor.GetX())

		out.SetLine(string(i.Lines[len(i.Lines)-1]))
		out.Devices["main"].(Display).RenderLine()
	case "NEWL":
		if Term.Cursor.GetY() == (Term.Lines-(Term.Lines/3))-1 {
			Term.Cursor.MoveTo(0, len(i.Lines)-1)
			Term.Cursor.SetX(0)

			i.AddLines(1)

			i.done <- struct{}{}
			return
		}

		Term.Cursor.MoveTo(0, len(i.Lines))
		Term.Cursor.AddY(len(i.Lines) - Term.Cursor.GetY())
		Term.Cursor.SetX(0)

		i.AddLines(1)

		i.done <- struct{}{}
		return
	case "TAB":
		if Term.Cursor.GetX() >= Term.Cols-8 ||
			Term.Cursor.GetX() < 0 {
			return
		}

		i.Lines[len(i.Lines)-1] = i.Lines[len(i.Lines)-1].Tab(Term.Cursor.GetX())
		Term.Cursor.AddX(4)

		out.SetLine(string(i.Lines[len(i.Lines)-1]))
		out.Devices["main"].(Display).RenderLine()
	}
}

// Movement buffer, an instruction buffer to hold arrow keystrokes
type mvbuf []byte

// Process for mvbuf. Actual implementation for the instructions
// in the mvbuf
func (mv mvbuf) process(i *Input) {
	switch string(mv) {
	case "UP":
		if Term.Cursor.GetY() <= 5 {
			return
		}

		if Term.Cursor.GetY() <= 5 && i.ScrollCount > 0 {
			i.ScrollBack()
			return
		}

		Term.Cursor.Up()
		Term.Cursor.AddY(-1)

		if Term.Cursor.GetX() > len(i.Lines[len(i.Lines)-1])+len(INPROMPT) {
			Term.Cursor.MoveTo(len(i.Lines[len(i.Lines)-1])+len(INPROMPT), Term.Cursor.GetY())
			Term.Cursor.SetX(len(i.Lines[len(i.Lines)-1]) + len(INPROMPT))
		}
	case "DOWN":
		if Term.Cursor.GetY() >= len(i.Lines)-1 {
			return
		}

		if Term.Cursor.GetY() == (Term.Lines/3)-1 && Term.Cursor.GetY() < len(i.Lines) && i.ScrollCount > 0 {
			i.Scroll()
			return
		}

		Term.Cursor.Down()
		Term.Cursor.AddY(1)

		if Term.Cursor.GetX() > len(i.Lines[len(i.Lines)-1])+len(INPROMPT) {
			Term.Cursor.MoveTo(len(i.Lines[len(i.Lines)-1])+len(INPROMPT), Term.Cursor.GetY())
			Term.Cursor.SetX(len(i.Lines[len(i.Lines)-1]) + len(INPROMPT))
		}
	case "RIGHT":
		if Term.Cursor.GetX() >= len(INPROMPT)+len(i.Lines[len(i.Lines)-1]) {
			return
		}

		Term.Cursor.Right()
		Term.Cursor.AddX(1)
	case "LEFT":
		if Term.Cursor.GetX() <= len(INPROMPT) {
			return
		}

		Term.Cursor.Left()
		Term.Cursor.AddX(-1)
	}
}

// Filter buffer, input is placed here to be filtered and processed
// into instructions for the input buffer system.
type fbuf []byte

// Process for fbuf.
func (f fbuf) process(i *Input) {
	f.filterInput(i)
}

// actual concurrent filtering of fbuf as necessary.
func (f fbuf) filterInput(i *Input) {
	var w sync.WaitGroup
	done := make(chan struct{}, 0)

	wg := &w
	wg.Add(6)

	go func() {
		go killCheck(i, wg)
		go DebugCheck(i, wg)
		go cmdCheck(i, wg)
		go parseArrows(i, wg)
		go otherSpecial(i, wg)
		go regularChars(i, wg)

		wg.Wait()
		done <- struct{}{}
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

func cmdCheck(i *Input, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x18" {
		wg.Done()
		return
	}

	if i.InCmdBar {
		i.InCmdBar = false
	} else {
		i.InCmdBar = true
	}

	// cmdBar := commandbar.NewCommandBar(3, Term.Cols-1, 1, Term.Lines-3, "black", "sharp")
	// cmdBar.Display()

	wg.Done()

}

// Filter through input byte for "Quit/Kill" KeyStroke: Ctrl-C.
func killCheck(i *Input, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x03" {
		wg.Done()
		return
	}

	Term.Normal()
	cli.Restore()
	os.Exit(0)
}

// Filter through input byte for "Debug" KeyStroke: Ctrl-B.
func DebugCheck(i *Input, wg *sync.WaitGroup) {
	if len(i.Fbuf) == 0 || string(i.Fbuf[0]) != "\x02" {
		wg.Done()
		return
	}

	if i.InDebug {
		i.Debugger.Off <- struct{}{}
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
	process(*Input)
}

// Process a slice of buffers and resetting their value back to
// empty, individually.
func ProccessBuffers(bufs []buffer, i *Input, wg *sync.WaitGroup) {

	for _, v := range bufs {
		v.process(i)

		reflect.ValueOf(v).Elem().SetBytes(
			[]byte(""),
		)
	}

	wg.Done()
	return
}
