package repl

import (
	"fmt"
	"os"

	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"github.com/BuprenorphineKid/golabs/pkg/commandbar"
)

const (
	WAIT = iota
	RESUME
	KILL
)

func ctrlB() {
	debug()
}

func debug() {

	term.Cursor.SavePos()

	term.Cursor.MoveTo(10, 15)
	fmt.Printf("%s", usr.Input.Lines)

	term.Cursor.MoveTo(10, 16)
	fmt.Printf("Lines = %d, Y = %d, X = %d",
		len(usr.Input.Lines),
		term.Cursor.Y,
		term.Cursor.X,
	)

	term.Cursor.MoveTo(10, 17)
	fmt.Printf("[4:] = %d", len(usr.Input.Lines[4:]))

	// select {
	// case pos := <-dbgCh:
	// 	term.Cursor.MoveTo(10, 17)
	// 	fmt.Printf("Lines[4:]: %d, POS: %d",
	// 		len(usr.Input.Lines[4:]),
	// 		pos,
	// 	)
	// default:
	// }

	term.Cursor.RestorePos()
}

func ctrlC() {
	kill()
}

func kill() {
	term.Normal()
	cli.Restore()
	os.Exit(0)
}

func ctrlX() {
	cmdbar()
}

func cmdbar() {
	term.Cursor.SavePos()
	defer term.Cursor.RestorePos()

	c := commandbar.NewCommandBar(3, term.Cols-1, 1, term.Lines-3, "black", "sharp")
	c.Display()
	cmd := c.Read()

	scrn.Win.Fill()
	scrn.Win.Draw()

	ExecuteCmd(usr, cmd)
}
