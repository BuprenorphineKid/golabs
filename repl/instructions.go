package repl

import "reflect"

type wbuf []byte

func (w wbuf) process(i *InOut) {
	i.write(w)
}

type rbuf []byte

func (r rbuf) process(i *InOut) {
	i.wbuf = append(i.wbuf, i.rbuf...)

}

type spbuf []byte

func (sp spbuf) process(i *InOut) {
	switch string(sp) {
	case "HOME":
		i.term.Cursor.Home()
	case "END":
		i.term.Cursor.End()
	case "BACK":
		i.line.Backspace()
	case "DEL":
		i.line.delChar()
	case "NEWL":
		i.wbuf = append(i.wbuf, "\x0a\x0d")
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
		if v == "" {
			return
		}

		v.process(in)

		switch v.(type) {
		case mvbuf:
			reflect.ValueOf(v.(mvbuf)).Elem().SetBytes(
				[]bytes(""),
			)
		case rbuf:
			reflect.ValueOf(v.(rbuf)).Elem().SetBytes(
				[]bytes(""),
			)
		case spbuf:
			reflect.ValueOf(v.(spbuf)).Elem().SetBytes(
				[]bytes(""),
			)
		}
	}
}
