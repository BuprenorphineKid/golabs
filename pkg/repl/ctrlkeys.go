package repl

import (
	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"github.com/BuprenorphineKid/golabs/pkg/commandbar"
	"os"
)

const (
	WAIT = iota
	RESUME
	KILL
)

func ctrlB() {

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
