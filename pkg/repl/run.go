package repl

import (
	"labs/pkg/cli"
	"labs/pkg/eval"
	"labs/pkg/readline"
	"labs/pkg/scripts"
	"labs/pkg/window"
	"sync"
)

// Singleton of Terminal.
var term = readline.Term

// Singleton of Frame.
var win = window.NewWindow(0, (term.Lines/3 + term.Lines/3), term.Lines/3, term.Cols, "grey", "thick")
var usr = NewUser(term)

// Instantiate objects, Start the main loop. This is the function
// you use to start the application.
func Run() {
	cli.Ready()
	defer cli.Restore()

	win.LoadScreen()

	term.Clear()
	term.RawMode()
	defer term.Normal()

	readline.Logo(usr.Input)

	readline.Init()

	scrn := window.NewScreen(win, term.Cursor)
	echo := readline.NewEcho(usr.Lab, usr.Input)

	win.Fill()
	win.Draw()

	scripter := scripts.NewHandler()
	scripter.Run()

	rCh := make(chan eval.Report)

	ctrlkeys := make(chan *sync.WaitGroup)

	go func() {
		for {
			ctrl := <-usr.Input.Ctrlkey
			var wg sync.WaitGroup
			wg.Add(1)

			ctrlkeys <- &wg

			switch ctrl {
			case "c":
				ctrlC()
			case "x":
				ctrlX()
			default:
			}

			wg.Done()
		}
	}()

	go func() {

		for {

			TakeInput(usr, echo, ctrlkeys)

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

			if len(scrn.Buffer) >= win.Height-2 {
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
func TakeInput(usr *User, echo *readline.Echo, ctrl chan *sync.WaitGroup) {
	if len(usr.Input.Lines) >= (term.Lines - (term.Lines / 3) - 1) {
		usr.Input.Scroll()

		input := readline.ReadLine(usr.Input, ctrl)

		DetermCmd(usr, string(*input), usr.FileLock)

		return
	}

	GiveOutput(echo)

	input := readline.ReadLine(usr.Input, ctrl)

	DetermCmd(usr, string(*input), usr.FileLock)
}
