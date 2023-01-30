package repl

import (
	"io"
	"os"
	"strings"

	"labs/pkg/cli"
)

// Application Input and Output meshed into a single construct that
// holds all the state relevent to working with the input and
// output.
type Input struct {
	reader   io.Reader
	writer   io.Writer
	Rbuf     rbuf
	Wbuf     wbuf
	Spbuf    spbuf
	Mvbuf    mvbuf
	Fbuf     fbuf
	lines    []line
	done     chan struct{}
	debugOff chan struct{}
	InDebug  bool
	Debugger *Debugger
}

// Constructor function for the Input Output struct.
func NewInput(t *cli.Terminal) *Input {
	i := Input{
		writer: os.Stdout,
		reader: os.Stdin,
	}

	i.Rbuf = rbuf("")
	i.Wbuf = wbuf("")
	i.Spbuf = spbuf("")
	i.Mvbuf = mvbuf("")
	i.Fbuf = fbuf("")
	i.lines = make([]line, 1, term.Lines)
	i.Debugger = new(Debugger)
	i.InDebug = false

	return &i
}

// InOut method for Adding n amount of lines.
func (i *Input) AddLines(n int) {
	newlines := make([]line, n, n)

	i.lines = append(i.lines, newlines...)
}

// The read method is used for recieving one byte of input at a
// time and appending it to the Filter buffer(Fbuf) unless
// it recieves an escape byte, in which case 2 more bytes
// will be read and then its entirety will be sent to the
// Filter buffer (Fbuf).
func (i *Input) read() {
	if term.IsRaw != true {
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

// The write method is used to write to the current line
// so that it can be processed by the output system.
func (i *Input) write(buf []byte) {
	if string(buf) == "" {
		return
	}

	i.lines[term.Cursor.Y] = i.lines[term.Cursor.Y].Insert(buf, term.Cursor.X)
}
