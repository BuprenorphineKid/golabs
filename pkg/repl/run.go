package repl

import (
	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"github.com/BuprenorphineKid/golabs/pkg/eval"
	"github.com/BuprenorphineKid/golabs/pkg/labs"
	"github.com/BuprenorphineKid/golabs/pkg/readline"
	"github.com/BuprenorphineKid/golabs/pkg/window"
	"os"
)

var home, _ = os.UserHomeDir()

// Singleton of Terminal.
var term = readline.Term

// Singleton of Frame.
var scrn = window.NewScreen(window.NewWindow(0, (term.Lines/3+term.Lines/3), term.Lines/3, term.Cols, "grey", "thick"), term.Cursor)
var usr = NewUser(term)

// Instantiate objects, Start the main loop. This is the function
// you use to start the application.
func Run() {
	cli.Ready()
	defer cli.Restore()

	term.RawMode()
	defer term.Normal()

	term.Clear()

	InitializeUI()

	rCh := make(chan eval.Report)

	go func() {
		for {
			ctrl := <-usr.Input.Ctrlkey

			usr.Input.CntrlCode <- WAIT

			switch ctrl {
			case "c":
				ctrlC()
			case "x":
				ctrlX()
			default:
			}

			usr.Input.CntrlCode <- RESUME

		}
	}()

	go func() {

		for {
			echo := readline.NewEcho(usr.Lab, usr.Input)

			TakeInput(echo)

			ev := eval.NewEvaluator(usr.Lab.Main)
			ev.Exec(rCh)

		}
	}()

	func() {
		for {
			info := <-rCh
			if !info.Ok {
				continue
			}

			scrn.Wrap(info.Results)

			if len(scrn.Buffer) >= scrn.Win.Height-2 {
				scrn.Scroll()
			}

			GiveOutput(scrn)
		}
	}()

}

// Shortcut to using Display.RenderLine() within the context
// of the Labs repl, that simply Takes in a *User and a
// string, modifies *User accordingly, and prints the
// passed in string after first printing the "output
// indication prompt."
//
// GiveOutput() is the counterpart to Take(. For each individual
// call to the Take() function, user does not hve to
// make a call to GiveOutput() before calling Take() again.
// Although it is recommended to try to.
func GiveOutput(out Outputter) {
	out.Display()
}

// Wraps GetLine() and handles the input accordingly.
// before returning a channel stream for Post-Evaluation
// printSlips.
//
// Mainly just to save you the hastle  of traversing my nooby
// unnecisarily complicated type relationships. But if you
// wish to handle it manually, you can.
//
// TakeInput() is to GetLine() what Give() is to
// Display.RenderLine().
func TakeInput(echo *readline.Echo) {
	if len(usr.Input.Lines) >= (term.Lines - (term.Lines / 3) - 1) {
		usr.Input.Scroll()

		input := readline.ReadLine(usr.Input)
		if input == nil {
			return
		}

		labs.DetermDecl(usr.Lab, string(*input.Line), usr.FileLock)

		return
	}

	GiveOutput(echo)

	input := readline.ReadLine(usr.Input)
	if input == nil {
		return
	}

	if input.Pos != len(usr.Input.Lines)-2 {
		usr.Lab.Replace(string(*input.Line), input.Pos)
		return
	}

	labs.DetermDecl(usr.Lab, string(*input.Line), usr.FileLock)
}
