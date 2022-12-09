package repl

import (
	"fmt"
	"labs/cli"
	"os"
	"reflect"
	"sync"
)

type wbuf []byte

func (w wbuf) process(i *InOut) {
	i.write(w)
}

type rbuf []byte

func (r rbuf) process(i *InOut) {
	i.Wbuf = append(i.Wbuf, i.Rbuf...)

}

type spbuf []byte

func (sp spbuf) process(i *InOut) {
	switch string(sp) {
	case "HOME":
		i.term.Cursor.Home()
	case "END":
		i.term.Cursor.End(len(i.lines[i.term.Cursor.Y]))
	case "BACK":
		i.lines[i.term.Cursor.Y].Backspace()
	case "DEL":
		i.lines[i.term.Cursor.Y].DelChar()
	case "NEWL":
		i.Wbuf = append(i.Wbuf, []byte("\x0a\x0d")...)
	}
}

type mvbuf []byte

func (mv mvbuf) process(i *InOut) {
	switch string(mv) {
	case "UP":
		i.term.Cursor.Up()
	case "DOWN":
		i.term.Cursor.Down()
	case "RIGHT":
		i.term.Cursor.Right()
	case "LEFT":
		i.term.Cursor.Left()
	}
}

type fbuf []byte

func (f fbuf) process(i *InOut) {
	f.filterInput(i)
}

func (f fbuf) filterInput(i *InOut) {
	var w sync.WaitGroup
	done := make(chan struct{}, 0)

	wg := &w
	wg.Add(4)

	go func() {
		go killCheck(i, wg)
		go testCheck(i, wg)
		go parseArrows(i, wg)
		go otherSpecial(i, wg)

		wg.Wait()
		done <- event{}
	}()

	select {
	case <-done:
		close(done)
	}
}

func otherSpecial(i *InOut, wg *sync.WaitGroup) {
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
	default:
		wg.Done()
		return
	}

	i.Rbuf = rbuf("")

	wg.Done()
}

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

func killCheck(i *InOut, wg *sync.WaitGroup) {
	if string(i.Fbuf[0]) != "\x03" {
		wg.Done()
		return
	}

	i.term.Normal()
	cli.Restore()
	os.Exit(3)
}

func testCheck(i *InOut, wg *sync.WaitGroup) {
	if string(i.Fbuf[0]) != "\x04" {
		wg.Done()
		return
	}

	fmt.Printf("BUFFERS\n_______\nrbuf: %s wbuf: %s mvbuf: %s spbuf: %s fbuf: %s\n",
		string(i.Rbuf),
		string(i.Wbuf),
		string(i.Mvbuf),
		string(i.Spbuf),
		string(i.Fbuf),
	)

	wg.Done()
}

type buffer interface {
	process(*InOut)
}

func ProccessBuffers(bufs []buffer, in *InOut) {
	for _, v := range bufs {
		v.process(in)

		reflect.ValueOf(v).Elem().SetBytes(
			[]byte(""),
		)
	}
}
