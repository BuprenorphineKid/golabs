package readline

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
	x := Term.Cols / 3
	y := Term.Lines / 3

	Term.Cursor.SavePos()
	Term.Cursor.Invisible()

	old := Term.Cursor.Y

	Term.Cursor.MoveTo(x+x, 1)
	Term.Cursor.CutRest()

	Term.Cursor.MoveTo(x+x, 2)
	Term.Cursor.CutRest()

	Term.Cursor.MoveTo(x+x, 3)
	Term.Cursor.CutRest()
	Term.Cursor.Y = 3

	for range i.Lines[5:] {
		Term.Cursor.CutRest()
		Term.Cursor.MoveTo(x+x, Term.Cursor.Y+1)
		Term.Cursor.AddY(1)
	}

	Term.Cursor.MoveTo(x+x, y-1)
	Term.Cursor.CutRest()

	for j := range reflect.VisibleFields(reflect.TypeOf(d.Stats).Elem()) {
		Term.Cursor.MoveTo(x+x, y+j)
		Term.Cursor.CutRest()
	}

	Term.Cursor.RestorePos()
	Term.Cursor.Y = old
	Term.Cursor.Normal()
}

func (d *Debugger) PosAndLines(i *Input) {
	x := Term.Cols / 3

	Term.Cursor.SavePos()
	Term.Cursor.Invisible()

	old := Term.Cursor.Y

	Term.Cursor.MoveTo(x+x, 1)
	Term.Cursor.CutRest()

	fmt.Printf("POS - |X: %d| |Y: %d|",
		Term.Cursor.X,
		Term.Cursor.Y,
	)

	Term.Cursor.MoveTo(x+x, 2)
	Term.Cursor.CutRest()

	fmt.Printf("LINES - |count: %d|",
		len(i.Lines),
	)

	Term.Cursor.MoveTo(x+x, 3)
	Term.Cursor.CutRest()
	Term.Cursor.Y = 3

	for n, v := range i.Lines[5:] {
		Term.Cursor.CutRest()
		fmt.Printf("%d - |%s|\n", n, string(v))
		Term.Cursor.MoveTo(x+x, Term.Cursor.Y+1)
		Term.Cursor.AddY(1)
		Term.Cursor.Normal()
	}

	Term.Cursor.RestorePos()
	Term.Cursor.Y = old
	Term.Cursor.Normal()

}

func (d *Debugger) MemStats(i *Input) {
	runtime.ReadMemStats(d.Stats)

	x := Term.Cols / 3
	y := Term.Lines / 3

	Term.Cursor.SavePos()
	Term.Cursor.Invisible()

	Term.Cursor.MoveTo(x+x, y-1)
	fmt.Print("|[MEMORY]|")

	for j, v := range reflect.VisibleFields(reflect.TypeOf(d.Stats).Elem()) {
		Term.Cursor.MoveTo(x+x, y+j)
		Term.Cursor.CutRest()
		fmt.Printf("|[%s: %d]|", v.Name, reflect.ValueOf(d.Stats).Elem().FieldByIndex(v.Index))
		if j == 22 {
			break
		}
	}

	Term.Cursor.RestorePos()
	Term.Cursor.Normal()
}
