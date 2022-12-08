package repl

import "reflect"

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

type buffer interface {
	process(*InOut)
}

func ProccessBuffers(bufs []buffer, in *InOut) {
	for _, v := range bufs {
		v.process(in)

		switch v.(type) {
		case mvbuf:
			if string(v.(mvbuf)) == "" {
				return
			}

			reflect.ValueOf(v.(*mvbuf)).Elem().SetBytes(
				[]byte(""),
			)
		case rbuf:
			if string(v.(rbuf)) == "" {
				return
			}

			reflect.ValueOf(v.(*rbuf)).Elem().SetBytes(
				[]byte(""),
			)
		case spbuf:
			if string(v.(spbuf)) == "" {
				return
			}

			reflect.ValueOf(v.(*spbuf)).Elem().SetBytes(
				[]byte(""),
			)
		case wbuf:
			if string(v.(wbuf)) == "" {
				return
			}

			reflect.ValueOf(v.(*wbuf)).Elem().SetBytes(
				[]byte(""),
			)
		}
	}
}
