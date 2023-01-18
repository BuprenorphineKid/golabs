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
	x := i.term.Cols / 3
	y := i.term.Lines / 3

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

	for range i.lines[5:] {
		i.term.Cursor.CutRest()
		i.term.Cursor.MoveTo(x+x, i.term.Cursor.Y+1)
		i.term.Cursor.AddY(1)
	}

	i.term.Cursor.MoveTo(x+x, y-1)
	i.term.Cursor.CutRest()

	for j := range reflect.VisibleFields(reflect.TypeOf(d.Stats).Elem()) {
		i.term.Cursor.MoveTo(x+x, y+j)
		i.term.Cursor.CutRest()
	}

	i.term.Cursor.RestorePos()
	i.term.Cursor.Y = old
	i.term.Cursor.Normal()
}

func (d *Debugger) PosAndLines(i *Input) {
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
		i.term.Cursor.Normal()
	}

	i.term.Cursor.RestorePos()
	i.term.Cursor.Y = old
	i.term.Cursor.Normal()

}

func (d *Debugger) MemStats(i *Input) {
	runtime.ReadMemStats(d.Stats)

	x := i.term.Cols / 3
	y := i.term.Lines / 3

	i.term.Cursor.SavePos()
	i.term.Cursor.Invisible()

	i.term.Cursor.MoveTo(x+x, y-1)
	fmt.Print("|[MEMORY]|")

	for j, v := range reflect.VisibleFields(reflect.TypeOf(d.Stats).Elem()) {
		i.term.Cursor.MoveTo(x+x, y+j)
		i.term.Cursor.CutRest()
		fmt.Printf("|[%s: %d]|", v.Name, reflect.ValueOf(d.Stats).Elem().FieldByIndex(v.Index))
		if j == 22 {
			break
		}
	}

	i.term.Cursor.RestorePos()
	i.term.Cursor.Normal()
}
