package repl

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

type Debugger struct {
	Off   chan struct{}
	Ready chan *sync.WaitGroup
	Stats *runtime.MemStats
}

func NewDebugger() *Debugger {
	d := Debugger{}
	d.Off = make(chan struct{})
	d.Ready = make(chan *sync.WaitGroup)
	d.Stats = new(runtime.MemStats)

	return &d
}

func (d *Debugger) DebugMode(i *Input) {
	d.PosAndLines(i)
	d.MemStats(i)
	for {
		select {
		case <-d.Off:
			d.CleanUp(i)
		case wg := <-d.Ready:
			d.PosAndLines(i)
			d.MemStats(i)
			wg.Done()
		}
	}
}

func (d *Debugger) CleanUp(i *Input) {
	x := term.Cols / 3
	y := term.Lines / 3

	term.Cursor.SavePos()
	term.Cursor.Invisible()

	old := term.Cursor.Y

	term.Cursor.MoveTo(x+x, 1)
	term.Cursor.CutRest()

	term.Cursor.MoveTo(x+x, 2)
	term.Cursor.CutRest()

	term.Cursor.MoveTo(x+x, 3)
	term.Cursor.CutRest()
	term.Cursor.Y = 3

	for range i.lines[5:] {
		term.Cursor.CutRest()
		term.Cursor.MoveTo(x+x, term.Cursor.Y+1)
		term.Cursor.AddY(1)
	}

	term.Cursor.MoveTo(x+x, y-1)
	term.Cursor.CutRest()

	for j := range reflect.VisibleFields(reflect.TypeOf(d.Stats).Elem()) {
		term.Cursor.MoveTo(x+x, y+j)
		term.Cursor.CutRest()
	}

	term.Cursor.RestorePos()
	term.Cursor.Y = old
	term.Cursor.Normal()
}

func (d *Debugger) PosAndLines(i *Input) {
	x := term.Cols / 3

	term.Cursor.SavePos()
	term.Cursor.Invisible()

	old := term.Cursor.Y

	term.Cursor.MoveTo(x+x, 1)
	term.Cursor.CutRest()

	fmt.Printf("POS - |X: %d| |Y: %d|",
		term.Cursor.X,
		term.Cursor.Y,
	)

	term.Cursor.MoveTo(x+x, 2)
	term.Cursor.CutRest()

	fmt.Printf("LINES - |count: %d|",
		len(i.lines),
	)

	term.Cursor.MoveTo(x+x, 3)
	term.Cursor.CutRest()
	term.Cursor.Y = 3

	for n, v := range i.lines[5:] {
		term.Cursor.CutRest()
		fmt.Printf("%d - |%s|\n", n, string(v))
		term.Cursor.MoveTo(x+x, term.Cursor.Y+1)
		term.Cursor.AddY(1)
		term.Cursor.Normal()
	}

	term.Cursor.RestorePos()
	term.Cursor.Y = old
	term.Cursor.Normal()

}

func (d *Debugger) MemStats(i *Input) {
	runtime.ReadMemStats(d.Stats)

	x := term.Cols / 3
	y := term.Lines / 3

	term.Cursor.SavePos()
	term.Cursor.Invisible()

	term.Cursor.MoveTo(x+x, y-1)
	fmt.Print("|[MEMORY]|")

	for j, v := range reflect.VisibleFields(reflect.TypeOf(d.Stats).Elem()) {
		term.Cursor.MoveTo(x+x, y+j)
		term.Cursor.CutRest()
		fmt.Printf("|[%s: %d]|", v.Name, reflect.ValueOf(d.Stats).Elem().FieldByIndex(v.Index))
		if j == 22 {
			break
		}
	}

	term.Cursor.RestorePos()
	term.Cursor.Normal()
}
