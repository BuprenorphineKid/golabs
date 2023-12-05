package repl

import (
	"fmt"
	"labs/pkg/syntax"
	"os"
	"os/exec"
)

type frame struct {
	x      int
	y      int
	height int
	width  int
	style  string
}

func NewFrame(x, y, h, w int, s string) *frame {
	var f frame
	f.x = x
	f.y = y
	f.height = h
	f.width = w
	f.style = s

	return &f
}

func (f frame) Draw() {
	term.Cursor.SavePos()

	var scrnLine string

	for i := 0; i < frm.width; i++ {
		scrnLine = scrnLine + " "
	}

	for i := 0; i < frm.height; i++ {
		term.Cursor.MoveTo(frm.x, frm.y+i)
		fmt.Print(syntax.OnGrey(scrnLine))
	}

	term.Cursor.RestorePos()

	cmd := exec.Command("tbox", fmt.Sprint(f.height), fmt.Sprint(f.width), fmt.Sprint(f.x), fmt.Sprint(f.y), f.style)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Run()
}
