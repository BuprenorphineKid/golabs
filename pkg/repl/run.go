package repl

import (
	"fmt"
	"labs/pkg/cli"
	"labs/pkg/eval"
	"labs/pkg/readline"
	"labs/pkg/scripts"
	"labs/pkg/syntax"
	"labs/pkg/window"
	"strings"
)

// Singleton of Terminal.
var term = readline.Term

// Singleton of Frame.
var win = window.NewWindow(0, (term.Lines/3 + term.Lines/3), term.Lines/3, term.Cols, "thick")

// Instantiate objects, Start the main loop. This is the function
// you use to start the application.
func Run() {
	cli.Ready()
	defer cli.Restore()

	win.LoadScreen()

	term.Clear()
	term.RawMode()
	defer term.Normal()

	usr := NewUser(term)
	readline.Logo(usr.Input)

	readline.Out.Register("main", readline.NewScreen())

	win.Fill()
	win.Draw()

	scripter := scripts.NewHandler()
	scripter.Run()

	rCh := make(chan eval.Report)

	go func() {
		for {
			TakeInput(usr)

			ev := eval.NewEvaluator(usr.Lab.Main)
			go ev.Exec(rCh)

		}
	}()

	func() {
		for {
			info := <-rCh
			if !info.Ok {
				continue
			}

			GiveOutput(usr, info.Results)
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
func GiveOutput(u *User, s string) {
	term.Cursor.SavePos()
	defer term.Cursor.RestorePos()

	win.Fill()
	defer win.Draw()

	ycurs := win.Y + 1
	xcurs := win.X + 3
	term.Cursor.MoveTo(xcurs, ycurs)

	if strings.Contains(s, "\n") {

		out := strings.Split(s, "\n")

		for _, v := range out {
			w := win.Width

			if len(v) > w {
				fmt.Print(syntax.OnGrey(syntax.White(v[:w])))

				ycurs++
				term.Cursor.MoveTo(xcurs, ycurs)

				fmt.Print(syntax.OnGrey(syntax.White(v[w:])))

				ycurs++
				term.Cursor.MoveTo(xcurs, ycurs)

				continue
			}

			fmt.Print(syntax.OnGrey(syntax.White(v)))
			return
		}
	}

	fmt.Print(syntax.OnGrey(syntax.White(s)))

	// output.devices["main"].(Display).RenderLine()

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
func TakeInput(usr *User) {
	if len(usr.Input.Lines) >= (term.Lines - (term.Lines / 3) - 1) {
		usr.Input.Scroll()

		input := readline.ReadLine(usr.Input)

		DetermCmd(usr, string(*input), usr.FileLock)

		return
	}

	if usr.Lab.InBody {
		readline.Out.Devices["main"].(readline.Display).PrintAndPrompt(&usr.Input.Lines, usr.Lab.Depth)
	} else {
		readline.Out.Devices["main"].(readline.Display).PrintInPrompt()
	}

	input := readline.ReadLine(usr.Input)

	DetermCmd(usr, string(*input), usr.FileLock)
}
