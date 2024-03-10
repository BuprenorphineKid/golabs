package window

import (
	"fmt"
	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"github.com/BuprenorphineKid/golabs/pkg/syntax"
)

type Window struct {
	Frame
	Color    string
	scrnLine string
}

func NewWindow(x, y, h, w int, c, s string) *Window {
	var win Window
	win.X = x
	win.Y = y
	win.Height = h
	win.Width = w
	win.Color = c
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

		switch w.Color {
		case "white":
			fmt.Print(syntax.OnWhite(w.scrnLine))
		case "black":
			fmt.Print(syntax.OnBlack(w.scrnLine))
		case "grey":
			fmt.Print(syntax.OnGrey(w.scrnLine))
		case "red":
			fmt.Print(syntax.OnRed(w.scrnLine))
		case "blue":
			fmt.Print(syntax.OnGrey(w.scrnLine))
		case "green":
			fmt.Print(syntax.OnGreen(w.scrnLine))
		case "yellow":
			fmt.Print(syntax.OnYellow(w.scrnLine))
		case "magenta":
			fmt.Print(syntax.OnMagenta(w.scrnLine))
		case "cyan":
			fmt.Print(syntax.OnCyan(w.scrnLine))
		}

	}
}
