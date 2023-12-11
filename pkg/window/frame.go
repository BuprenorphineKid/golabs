package window

import (
	"fmt"
	"os"
	"os/exec"
)

type Frame struct {
	X      int
	Y      int
	Height int
	Width  int
	Style  string
}

func NewFrame(x, y, h, w int, s string) *Frame {
	var f Frame
	f.X = x
	f.Y = y
	f.Height = h
	f.Width = w
	f.Style = s

	return &f
}

func (f Frame) Draw() {
	cmd := exec.Command("tbox", fmt.Sprint(f.Height), fmt.Sprint(f.Width), fmt.Sprint(f.X), fmt.Sprint(f.Y), f.Style)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Run()
}
