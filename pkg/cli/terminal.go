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
	Cursor   Cursor
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

	c := exec.Command("tput", "cols")
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
