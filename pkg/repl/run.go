package repl

import (
	"labs/pkg/cli"
	"labs/pkg/labs"
	"labs/pkg/scripts"
	"sync"
)

var term = cli.NewTerminal()

// User Struct for keeping count of Hist CmdCount, yada yada.
// User holds all the main ingredients to run the show that
// are exvlusive to a user. Think env, files, input/output etc.
type User struct {
	Name     string
	Input    *Input
	Lab      *labs.Lab
	Logger   *Log
	done     chan struct{}
	FileLock sync.Locker
}

// Creates a new User object and returns a pointer to it.
func NewUser(t *cli.Terminal) *User {

	var u User

	u.Input = NewInput(t)
	u.Lab = labs.NewLab()

	u.FileLock = new(sync.Mutex)

	return &u
}

// Set the Users Name.
func (u *User) setName(name string) {
	u.Name = name
}

// Instantiate objects, Start the main loop, and start up the
// Evaluation Controller. This is the function you use to
// start the application
func Run() {
	cli.Ready()
	defer cli.Restore()

	term.Clear()
	term.RawMode()
	defer term.Normal()

	lgr := NewLog()
	usr := NewUser(term)
	logo(usr.Input)

	usr.Logger = lgr

	output.Register("main", newScreen())

	scripter := scripts.NewHandler()
	scripter.Run()

	bash := scripts.NewLanguage("bash")

	for {
		outCh := Take(usr)

		select {
		case vfy := <-outCh:
			func() {
				if !vfy.ok {
					return
				}

				Give(usr, vfy.results)
			}()
		}

		scripter.Do <- scripts.Exec(bash, "scripts/remove_print.sh")
	}
}

// Shortcut to using Display.RenderLine() within the context
// of the Labs repl, that simply Takes in a *User and a
// string, modifies *User accordingly, and prints the
// passed in string after first printing the "output
// indication prompt."
//
// Give() is the counterpart to Take(. For each individual
// call to the Take() function, user does not hve to
// make a call to Give() before calling Take() again.
// Although it is recommended to try to.
func Give(u *User, s string) {
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
// Take() is to GetLine() what Give() is to
// Display.RenderLine().
func Take(usr *User) chan printSlip {
	if usr.Lab.InBody {
		output.devices["main"].(Display).PrintAndPrompt(&usr.Input.lines, usr.Lab.Depth)
	} else {
		output.devices["main"].(Display).PrintInPrompt()
	}

	input := GetLine(usr.Input)

	DetermCmd(usr, string(*input), usr.FileLock)

	output := make(chan printSlip)
	e := NewEvaluator(usr.Lab.Main)

	go e.Exec(output)

	return output
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
