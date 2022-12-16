package cli

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Just a standard terminal representation for storing cursor position and
// a few other env vars.
type Terminal struct {
	IsRaw    bool
	Cursor   cursor
	OldState *term.State
	Cols     int
	Lines    int
}

func NewTerminal() *Terminal {
	t := Terminal{}
	t.IsRaw = false
	t.Cols, t.Lines = size()
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
	c.Stdin = os.Stdin

	strC, err := c.Output()
	if err != nil {
		panic(err)
	}
	str[0] = strings.TrimSpace(string(strC))
	size64[0], _ = strconv.ParseInt(str[0], 0, 0)
	size[0] = int(size64[0])

	l := exec.Command("tput", "lines")
	l.Stdin = os.Stdin

	var strL []byte
	strL, err = l.Output()
	if err != nil {
		panic(err)
	}
	str[1] = strings.TrimSpace(string(strL))
	size64[1], _ = strconv.ParseInt(str[1], 0, 0)
	size[1] = int(size64[1])

	return size[0], size[1]
}

type cursor struct {
	X int
	Y int
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
	c.X = 0
	c.Y = 0
}

func (c *cursor) Home(prompt int) {
	print("\033[", c.Y, ";", prompt, "H")
}

func (c *cursor) End(end int) {
	print("\033[", c.Y, ";", end, "H")
}

func (c *cursor) Left() {
	print("\033[D")
}

func (c *cursor) Right() {
	print("\033[C")
}

func (c *cursor) Up() {
	print("\033[A")
}

func (c *cursor) Down() {
	print("\033[B")
}

func (c *cursor) AddX(n int) {
	c.X = c.X + n
}

func (c *cursor) AddY(n int) {
	c.Y = c.Y + n
}

func (c *cursor) MoveTo(x int, y int) {
	print("\033[", y, ";", x, "H")
}

func (c *cursor) CutRest() {
	print("\033[0K")
}

func (c *cursor) CutFirst() {
	print("\033[1K")
}

func (c *cursor) CutLine() {
	print("\033[2K")
}

func (c *cursor) Invisible() {
	print("\033[?25l")
}

func (c *cursor) Normal() {
	print("\033[?12l\033[?25h")
}
