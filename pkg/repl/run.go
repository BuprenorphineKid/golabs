package repl

import (
	"labs/pkg/cli"
	"labs/pkg/scripts"
	"labs/pkg/syntax"
	"sync"
)

// Singleton of Terminal.
var term = cli.NewTerminal()

// Instantiate objects, Start the main loop. This is the function
// you use to start the application.
func Run() {
	cli.Ready()
	defer cli.Restore()

	term.Clear()
	term.RawMode()
	defer term.Normal()

	usr := NewUser(term)
	logo(usr.Input)

	output.Register("main", newScreen())

	go func() {
		scripter := scripts.NewHandler()
		scripter.Run()

		for {
			TakeInput(usr)

			eval := NewEvaluator(usr.Lab.Main)
			rCh := make(chan report)
			go eval.Exec(rCh)

			select {
			case info := <-rCh:
				if !info.ok {
					break
				}

				GiveOutput(usr, info.results)

				scripter.Do <- scripts.Exec(scripts.NewLanguage("bash"), "scripts/bash/extract_vars.sh")
			}
		}
	}()

	func() {
		scripter := scripts.NewHandler()
		scripter.Run()

		for {
			if usr.Lab.History.Last == nil {
				continue
			}

			if syntax.IsFuncCall(*usr.Lab.History.Last) {
				scripter.Do <- scripts.Exec(scripts.NewLanguage("bash"), "scripts/bash/clear_main.sh")

			}
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
	output.SetLine(s)
	output.devices["main"].(Display).PrintOutPrompt()
	output.devices["main"].(Display).RenderLine()

	u.Input.AddLines(2)
	term.Cursor.AddY(2)
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
	if usr.Lab.InBody {
		output.devices["main"].(Display).PrintAndPrompt(&usr.Input.lines, usr.Lab.Depth)
	} else {
		output.devices["main"].(Display).PrintInPrompt()
	}

	input := GetLine(usr.Input)

	DetermCmd(usr, string(*input), usr.FileLock)
}

// Most users will be more inclined to use Take() instead
// because it does the dirty work for you.
//
// Start Input Loop that after its done reading input and
// filling buffers, concurrently processes each.
func GetLine(i *Input) *line {
	for {
		if i.InDebug {
			var dbwg sync.WaitGroup
			dbwg.Add(1)
			i.Debugger.Ready <- &dbwg
			dbwg.Wait()
		}

		i.done = make(chan struct{}, 1)
		i.read()

		var nlwg sync.WaitGroup
		nlwg.Add(1)

		var bufs = []buffer{&i.Fbuf, &i.Rbuf, &i.Spbuf, &i.Mvbuf, &i.Wbuf}
		go ProccessBuffers(bufs, i, &nlwg)

		nlwg.Wait()

		output.SetLine(string(i.lines[term.Cursor.Y]))
		output.devices["main"].(Display).RenderLine()

		select {
		case <-i.done:
			return &i.lines[term.Cursor.Y-1]
		default:
			continue
		}

	}
}
