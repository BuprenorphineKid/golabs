package repl

import (
	"io"
	"os"
	"strings"

	"labs/cli"
)

type event struct{}

type InOut struct {
	reader io.Reader
	writer io.Writer
	rbuf   []byte
	buf    [1]byte
	wbuf   []byte
	term   *cli.Terminal
	lines  []string
}

func newInOut(t *cli.Terminal) *InOut {
	i := InOut{
		writer: os.Stdout,
		reader: os.Stdin,
	}
	i.term = t
	i.rbuf = make([]byte, 0)
	i.wbuf = make([]byte, 0)
	i.lines = make([]string, i.term.Lines, i.term.Lines)

	return &i
}

func (i *InOut) read() []byte {
	i.rbuf = make([]byte, 0)

	if i.term.IsRaw != true {
		return i.rbuf

	}
	done := make(chan struct{}, 2)

	for {
		_, err := i.reader.(*os.File).Read(i.buf[:])
		if err != nil {
			panic(err)
		}

		i.rbuf = append(i.rbuf, i.buf[0])

		go i.killCheck(done)
		go i.specialCheck(done)

		j := 2
		for {
			select {
			case <-done:
				j--
			default:
			}
			if j == 0 {
				break
				close(done)
			}
		}

		if len(i.rbuf) > 0 {
			break
		}
	}
	return i.rbuf
}

func (i *InOut) killCheck(ch chan struct{}) {
	if string(i.rbuf) != "\x03" {
		ch <- event{}
		return
	}

	i.term.Normal()
	cli.Restore()
	os.Exit(3)
}

func (i *InOut) specialCheck(ch chan struct{}) {
	switch {
	case string(i.rbuf) == "\x0a" || string(i.rbuf) == "\x0d": //newline/carriage-return
		i.rbuf = []byte(strings.Replace(string(i.rbuf), string(i.rbuf), "\x0a\x0d", 1))
		i.term.Cursor.AddY(1)
	case string(i.rbuf) == "\x7F": //backspace
		i.rbuf = []byte(strings.Replace(string(i.rbuf), string(i.rbuf), "\x08 \x08", 1))
		i.term.Cursor.AddX(-1)
	case string(i.rbuf) == "\033[H": //home
		i.rbuf = []byte("")
		i.term.Cursor.Home()
		i.term.Cursor.X = 0
	case string(i.rbuf) == "\033[F": //end
		i.rbuf = []byte("")
		i.term.Cursor.End(i.term.Cols)
		i.term.Cursor.X = i.term.Cols
	case string(i.rbuf) == "\033[A": //up
		i.term.Cursor.AddY(-1)
		i.rbuf = []byte("")
	case string(i.rbuf) == "\033[B": //down
		i.term.Cursor.AddY(1)
		i.rbuf = []byte("")
	case string(i.rbuf) == "\033[C": //right
		i.term.Cursor.AddX(1)
		i.rbuf = []byte("")
	case string(i.rbuf) == "\x1b\x5b\x44": //left
		i.term.Cursor.AddX(-1)
		i.rbuf = []byte("")
	default:
		i.term.Cursor.X += len(i.rbuf)
	}

	ch <- event{}
}

func (i *InOut) write(buf []byte) {
	_, err := i.writer.(*os.File).Write(buf)

	if err != nil {
		panic(err)
	}

	func() {
		i.lines[i.term.Cursor.Y] += string(buf)
	}()
}

func StartInputLoop(inout *InOut) string {
	printLineLogo()

	inout.term.RawMode()
	for {
		inout.wbuf = inout.read()
		go inout.write(inout.wbuf)

		if string(inout.wbuf) == "\x0a\x0d" {
			return inout.lines[inout.term.Cursor.Y]
		}
	}
}
