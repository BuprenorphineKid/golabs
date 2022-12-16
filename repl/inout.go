package repl

import (
	"io"
	"os"
	"strings"
	"sync"

	"labs/cli"
)

// Application Input and Output meshed into a single construct that
// holds all the state relevent to working with the input and
// output.
type InOut struct {
	reader   io.Reader
	writer   io.Writer
	Rbuf     rbuf
	Wbuf     wbuf
	Spbuf    spbuf
	Mvbuf    mvbuf
	Fbuf     fbuf
	term     *cli.Terminal
	lines    []line
	done     chan struct{}
	debugOff chan struct{}
	InDebug  bool
	Debugger *Debugger
}

// Constructor function for the Input Output struct.
func NewInOut(t *cli.Terminal) *InOut {
	i := InOut{
		writer: os.Stdout,
		reader: os.Stdin,
	}

	i.term = t
	i.Rbuf = rbuf("")
	i.Wbuf = wbuf("")
	i.Spbuf = spbuf("")
	i.Mvbuf = mvbuf("")
	i.Fbuf = fbuf("")
	i.lines = make([]line, 1, i.term.Lines)
	i.Debugger = new(Debugger)
	i.InDebug = false

	return &i
}

// InOut method for Adding n amount of lines.
func (i *InOut) AddLines(n int) {
	newlines := make([]line, n, n)

	i.lines = append(i.lines, newlines...)
}

// the read method is used for recieving one byte of input at a
// time and appending it to the Filter buffer(Fbuf) unless
// it recieves an escape byte, in which case 2 more bytes
// will be read and then its entirety will be sent to the
// Filter buffer (Fbuf)
func (i *InOut) read() {
	if i.term.IsRaw != true {
		panic("Not able to enter into raw mode :(.")
	}

	var buf [1]byte

	_, err := i.reader.(*os.File).Read(buf[:])

	if err != nil {
		panic(err)
	}

	i.Fbuf = append(i.Fbuf, buf[:]...)

	for strings.HasPrefix(string(i.Fbuf), "\x1b") && len(i.Fbuf) < 3 {
		var buf [1]byte

		_, err := i.reader.(*os.File).Read(buf[:])

		if err != nil {
			panic(err)
		}

		i.Fbuf = append(i.Fbuf, buf[:]...)
	}

	if strings.HasPrefix(string(i.Fbuf), "~") {
		i.Fbuf = fbuf(strings.Trim(string(i.Fbuf), "~"))
	}
}

// the write method is used to simultaneously write to the screen
// and to the current line.
func (i *InOut) write(buf []byte) {
	if string(buf) == "" {
		return
	}

	i.lines[i.term.Cursor.Y] = i.lines[i.term.Cursor.Y].Insert(buf, i.term.Cursor.X)

	Refresh(i)
}

// Start Input Loop that after its done reading input and
// filling buffers, concurrently processes each.
func StartInputLoop(i *InOut) line {
	printLineLogo(i)

	i.term.Cursor.X = len(LINELOGO)

	for {
		if i.InDebug {
			var dbwg sync.WaitGroup
			dbwg.Add(1)

			i.Debugger.Ready <- &dbwg

			dbwg.Wait()
		}

		i.done = EventChan(1)

		i.read()

		var nlwg sync.WaitGroup
		nlwg.Add(1)

		var bufs = []buffer{&i.Fbuf, &i.Rbuf, &i.Spbuf, &i.Mvbuf, &i.Wbuf}
		go ProccessBuffers(bufs, i, &nlwg)

		nlwg.Wait()

		select {
		case <-i.done:
			return i.lines[i.term.Cursor.Y-1]
		default:
			continue
		}

	}
}
