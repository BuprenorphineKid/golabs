package commandbar

import (
	"fmt"
	"github.com/BuprenorphineKid/golabs/pkg/window"
	"os"
	"os/exec"
)

type bar struct {
	h, w, x, y int
	color      string
	style      string
	window.Window
}

func newBar(h, w, x, y int, c, s string) bar {
	var b bar
	b.h = h
	b.w = w
	b.x = x
	b.y = y
	b.style = s
	b.color = c
	b.Window = *window.NewWindow(x, y, h, w, c, s)

	return b
}

func (b bar) display() {
	b.LoadScreen()
	b.Fill()
	cmd := exec.Command("tbox", fmt.Sprint(b.h), fmt.Sprint(b.w), fmt.Sprint(b.x), fmt.Sprint(b.y), b.style)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Run()
}
