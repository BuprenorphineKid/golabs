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
}

func newInOut(t *cli.Terminal) *InOut {
	i := InOut{
		writer: os.Stdout,
		reader: os.Stdin,
	}
	i.rbuf = make([]byte, 0)
	i.wbuf = make([]byte, 0)
	i.term = t

	return &i
}

func (i *InOut) read() []byte {
	i.rbuf = make([]byte, 0)

	if i.term.IsRaw == true {
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
	case string(i.rbuf) == "\x0a" || string(i.rbuf) == "\x0d":
		i.rbuf = []byte(strings.Replace(string(i.rbuf), string(i.rbuf), "\x0a\x0d", 1))
		ch <- event{}
	case string(i.rbuf) == "\x7F":
		i.rbuf = []byte(strings.Replace(string(i.rbuf), string(i.rbuf), "\x08 \x08", 1))
		ch <- event{}
	case string(i.rbuf) == "\033[H":

	default:
		ch <- event{}
	}
}

func (i *InOut) write(buf []byte) {
	_, err := i.writer.(*os.File).Write(buf)

	if err != nil {
		panic(err)
	}
}

func StartInputLoop(inout *InOut) {
	inout.term.Clear()

	welcome()
	printLineLogo()

	inout.term.RawMode()
	for {
		inout.wbuf = inout.read()
		inout.write(inout.wbuf)

		if string(inout.wbuf) == "\n\r" {
			printLineLogo()
		}
	}
}
