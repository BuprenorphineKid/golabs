package repl

import (
	"io"
	"os"
	"strings"
	"sync"

	"labs/cli"
)

type InOut struct {
	reader io.Reader
	writer io.Writer
	Rbuf   rbuf
	Wbuf   wbuf
	Spbuf  spbuf
	Mvbuf  mvbuf
	Fbuf   fbuf
	term   *cli.Terminal
	lines  []line
	done   chan struct{}
}

func newInOut(t *cli.Terminal) *InOut {
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

	return &i
}

func (i *InOut) AddLines(n int) {
	newlines := make([]line, n, n)

	i.lines = append(i.lines, newlines...)
}

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

func (i *InOut) write(buf []byte) {
	if string(buf) == "" {
		return
	}

	_, err := i.writer.(*os.File).Write(buf)

	if err != nil {
		panic(err)
	}

	go func() {
		i.lines[i.term.Cursor.Y] += line(buf)
	}()
}

func StartInputLoop(i *InOut) *line {
	if i.term.Cursor.X < len(LINELOGO) {
		printLineLogo(i)
	}

	for {
		i.done = EventChan(1)

		i.read()

		var wg sync.WaitGroup
		wg.Add(1)

		var bufs = []buffer{&i.Fbuf, &i.Rbuf, &i.Spbuf, &i.Mvbuf, &i.Wbuf}
		go ProccessBuffers(bufs, i, &wg)

		wg.Wait()

		select {
		case <-i.done:
			return &i.lines[i.term.Cursor.Y]
		default:
			continue
		}

	}
}
