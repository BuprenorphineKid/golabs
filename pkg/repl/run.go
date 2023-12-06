package repl

import (
	"fmt"
	"labs/pkg/cli"
	"labs/pkg/scripts"
	"labs/pkg/syntax"
	"strings"
	"sync"
)

// Singleton of Terminal.
var term = cli.NewTerminal()

// Singleton of Frame.
var frm = NewFrame(0, ((term.Lines / 3) + (term.Lines / 3)), term.Lines/3, term.Cols, "thick")

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

	frm.Fill()
	frm.Draw()

	scripter := scripts.NewHandler()
	scripter.Run()

	rCh := make(chan report)

	go func() {
		for {
			TakeInput(usr)

			eval := NewEvaluator(usr.Lab.Main)
			go eval.Exec(rCh)

		}
	}()

	func() {
		for {
			info := <-rCh
			if !info.ok {
				continue
			}

			go GiveOutput(usr, info.results)
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
func GiveOutput(u *User, str string) {
	s := syntax.OnGrey(syntax.White(str))

	func() {
		term.Cursor.SavePos()

		for l := frm.y + 1; l < term.Lines; l++ {
			term.Cursor.MoveTo(frm.x+3, l)
			term.Cursor.CutRest()

			frm.Draw()
			frm.Fill()
		}

		term.Cursor.RestorePos()
	}()

	term.Cursor.SavePos()
	defer term.Cursor.RestorePos()
	defer frm.Draw()

	ycurs := frm.y + 1
	term.Cursor.MoveTo(frm.x+1, ycurs)

	if strings.Contains(s, "\n") {

		out := strings.Split(s, "\n")

		for _, v := range out {
			w := frm.width - 2
			if len(v) >= w {
				fmt.Print(v[:w])

				ycurs++
				term.Cursor.MoveTo(frm.x+1, ycurs)
				ycurs++

				fmt.Print(v[w:])
			}
			fmt.Print(v)
		}
	}

	fmt.Print(s)

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
