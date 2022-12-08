package repl

import (
	"io"
	"os"
	"sync"

	"labs/cli"
)

type event struct{}

type InOut struct {
	reader io.Reader
	writer io.Writer
	Rbuf   rbuf
	Wbuf   wbuf
	Spbuf  spbuf
	Mvbuf  mvbuf
	term   *cli.Terminal
	lines  []line
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
	i.lines = make([]line, 1, i.term.Lines)

	return &i
}

func (i *InOut) AddLine(n int) {
	newlines := make([]line, n, n)

	i.lines = append(i.lines, newlines...)
}

func (i *InOut) read() {
	if i.term.IsRaw != true {
		panic("Not able to enter into raw mode :(.")
	}

	var buf [1]byte
	//	var checkbuf []byte

	done := make(chan struct{}, 0)

	for {
		_, err := i.reader.(*os.File).Read(buf[:])

		if err != nil {
			panic(err)
		}

		//		checkbuf = append(checkbuf, buf...)

		var w sync.WaitGroup
		wg := &w
		wg.Add(3)

		go func() {
			go i.killCheck(buf[:], wg)
			go i.parseArrows(buf[:], wg)
			go i.otherSpecial(buf[:], wg)

			wg.Wait()
			done <- event{}
		}()

		select {
		case <-done:
			close(done)
		}

		i.Rbuf = append(i.Rbuf, buf[:]...)

		if len(i.Rbuf) > 0 {
			break
		}
	}
	return
}

func (i *InOut) killCheck(buf []byte, wg *sync.WaitGroup) {
	if string(buf[0]) != "\x03" {
		wg.Done()
		return
	}

	i.term.Normal()
	cli.Restore()
	os.Exit(3)
}

func (i *InOut) otherSpecial(buf []byte, wg *sync.WaitGroup) {
	switch string(buf[0]) {
	case "\033":
		switch string(buf[1]) {
		case "[":
			switch string(buf[2]) {
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
	}

	wg.Done()
}

func (i *InOut) parseArrows(buf []byte, wg *sync.WaitGroup) {
	if string(buf[0]) != "\x1b" {
		wg.Done()
		return
	}

	switch string(buf[1]) {
	case "[":
		switch string(buf[2]) {
		case "A":
			i.Mvbuf = mvbuf("UP")
		case "B":
			i.Mvbuf = mvbuf("DOWN")
		case "C":
			i.Mvbuf = mvbuf("RIGHT")
		case "D":
			i.Mvbuf = mvbuf("LEFT")
		}
	}

	wg.Done()
}

func (i *InOut) write(buf []byte) {
	_, err := i.writer.(*os.File).Write(buf)

	if err != nil {
		panic(err)
	}

	func() {
		i.lines[i.term.Cursor.Y] += line(buf)
	}()
}

func StartInputLoop(i *InOut) *line {
	printLineLogo()

	i.term.RawMode()
	for {
		i.read()

		var bufs = []buffer{&i.Rbuf, &i.Spbuf, &i.Mvbuf, &i.Wbuf}

		go ProccessBuffers(bufs, i)

		if string(i.Wbuf) == "\x0a\x0d" {
			return &i.lines[i.term.Cursor.Y]
		}
	}
}
