package window

import (
	"fmt"
	"labs/pkg/cli"
	"labs/pkg/syntax"
)

type Window struct {
	Frame
	scrnLine string
}

func NewWindow(x, y, h, w int, s string) *Window {
	var win Window
	win.X = x
	win.Y = y
	win.Height = h
	win.Width = w
	win.Style = s

	return &win
}

func (w *Window) LoadScreen() {
	for i := 0; i < w.Width; i++ {
		w.scrnLine = w.scrnLine + " "
	}
}

func (w Window) Fill() {
	term := cli.NewTerminal()

	for i := 0; i < w.Height; i++ {
		term.Cursor.MoveTo(w.X, w.Y+i)
		fmt.Print(syntax.OnGrey(w.scrnLine))
	}
}
