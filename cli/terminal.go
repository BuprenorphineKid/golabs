package cli

import (
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/term"
)

// Just a standard terminal representation for storing cursor position and
// a few other env vars.
type Terminal struct {
	IsRaw    bool
	Cursor   cursor
	OldState *term.State
	cols     int
	lines    int
}

func NewTerminal() *Terminal {
	t := Terminal{}
	t.IsRaw = false
	t.cols, t.lines = size()
	return &t
}

func (t *Terminal) RawMode() {
	var err error

	t.OldState, err = term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	t.IsRaw = true
}

func (t *Terminal) Normal() {
	term.Restore(int(os.Stdin.Fd()), t.OldState)

	t.IsRaw = false
}

func (t *Terminal) Clear() {
	print("\033[H\033[2J\033[3J")

	t.Cursor.TrueHome()
}

func Ready() {
	print("\033[?1049h\033[22;0;0t")
}

func Restore() {
	print("\033[?1049l\033[23;0;0t")
}

func size() (int, int) {
	var size64 [2]int64
	var size [2]int
	var str [2]string

	c := exec.Command("tput", "columns")

	strC, err := c.Output()
	if err != nil {
		panic(err)
	}
	str[0] = string(strC)

	l := exec.Command("tput", "lines")

	var strL []byte
	strL, err = l.Output()
	if err != nil {
		panic(err)
	}
	str[1] = string(strL)

	for i := range size {
		size64[i], _ = strconv.ParseInt(str[i], 0, 0)
		size[i] = int(size64[i])
	}

	return size[0], size[1]
}

type cursor struct {
	x int
	y int
}

func newCursor() *cursor {
	c := cursor{}
	return &c
}

func (c *cursor) SavePos() {
	print("\0337")
}

func (c *cursor) RestorePos() {
	print("\0338")
}

func (c *cursor) TrueHome() {
	print("\033[0;0H")
	c.x = 0
	c.y = 0
}

func (c *cursor) Home() {
	print("\033[", c.y, ";0H")
}
