package repl

import (
	"fmt"
	"sync"
)

type Debugger struct {
	Off   chan struct{}
	Ready chan *sync.WaitGroup
}

func NewDebugger() *Debugger {
	d := Debugger{}
	d.Off = EventChan(1)
	d.Ready = make(chan *sync.WaitGroup)

	return &d
}

func (d *Debugger) DebugMode(i *InOut) {
	x := i.term.Cols / 3

	i.term.Cursor.SavePos()
	i.term.Cursor.Invisible()

	old := i.term.Cursor.Y

	i.term.Cursor.MoveTo(x+x, 1)
	i.term.Cursor.CutRest()

	fmt.Printf("POS - |X: %d| |Y: %d|",
		i.term.Cursor.X,
		i.term.Cursor.Y,
	)

	i.term.Cursor.MoveTo(x+x, 2)
	i.term.Cursor.CutRest()

	fmt.Printf("LINES - |count: %d|",
		len(i.lines),
	)

	i.term.Cursor.MoveTo(x+x, 3)
	i.term.Cursor.CutRest()
	i.term.Cursor.Y = 3

	for n, v := range i.lines[5:] {
		i.term.Cursor.CutRest()
		fmt.Printf("%d - |%s|\n", n, string(v))
		i.term.Cursor.MoveTo(x+x, i.term.Cursor.Y+1)
		i.term.Cursor.AddY(1)
	}

	i.term.Cursor.RestorePos()
	i.term.Cursor.Y = old
	i.term.Cursor.Normal()

	for {
		select {
		case <-d.Off:
			i.term.Cursor.SavePos()
			i.term.Cursor.Invisible()

			old := i.term.Cursor.Y

			i.term.Cursor.MoveTo(x+x, 1)
			i.term.Cursor.CutRest()

			i.term.Cursor.MoveTo(x+x, 2)
			i.term.Cursor.CutRest()

			i.term.Cursor.MoveTo(x+x, 3)
			i.term.Cursor.CutRest()
			i.term.Cursor.Y = 3

			for _ = range i.lines[5:] {
				i.term.Cursor.CutRest()
				i.term.Cursor.MoveTo(x+x, i.term.Cursor.Y+1)
				i.term.Cursor.AddY(1)
			}

			i.term.Cursor.RestorePos()
			i.term.Cursor.Y = old
			i.term.Cursor.Normal()

			return
		case wg := <-d.Ready:
			i.term.Cursor.SavePos()
			i.term.Cursor.Invisible()

			old := i.term.Cursor.Y

			i.term.Cursor.MoveTo(x+x, 1)
			i.term.Cursor.CutRest()

			fmt.Printf("POS - |X: %d| |Y: %d|",
				i.term.Cursor.X,
				i.term.Cursor.Y,
			)

			i.term.Cursor.MoveTo(x+x, 2)
			i.term.Cursor.CutRest()

			fmt.Printf("LINES - |count: %d|",
				len(i.lines),
			)

			i.term.Cursor.MoveTo(x+x, 3)
			i.term.Cursor.CutRest()
			i.term.Cursor.Y = 3

			for n, v := range i.lines[5:] {
				i.term.Cursor.CutRest()
				fmt.Printf("%d - |%s|\n", n, string(v))
				i.term.Cursor.MoveTo(x+x, i.term.Cursor.Y+1)
				i.term.Cursor.AddY(1)
			}

			i.term.Cursor.RestorePos()
			i.term.Cursor.Y = old
			i.term.Cursor.Normal()

			wg.Done()
		default:
		}
	}
}
