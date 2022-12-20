package repl

import (
	"labs/cli"
	"sync"
)

// User Struct for keeping count of Hist CmdCount, yada yada.
// User holds all the main ingredients to run the show that
// are exvlusive to a user. Think env, files, input/output etc.
type User struct {
	CmdCount  int
	Name      string
	CmdHist   []string
	InBody    bool
	NestDepth int
	Eval      *Eval
	InOut     *InOut
	Lab       *Lab
	Shader    Shader
}

// Creates a new User object and returns a pointer to it.
func NewUser(t *cli.Terminal) *User {
	h := make([]string, 1, 150)

	var u User

	u.CmdCount = 1
	u.CmdHist = h
	u.CmdHist[0] = "begin"

	u.InOut = NewInOut(t)
	u.Lab = NewLab()
	u.Eval = NewEval()
	u.Shader = newHiLiter(&u.InOut.lines)

	u.InBody = false
	u.NestDepth = 0

	return &u
}

// Set the Users Name.
func (u *User) setName(name string) {
	u.Name = name
}

// Increment CmdCount and add desired command to Hist.
func (u *User) addCmd(cmd string) {
	u.CmdHist = append(u.CmdHist, cmd)
	u.CmdCount++
}

// The main recursive loop at the tozp level.
func repl(usr *User) {
	input := StartInputLoop(usr)

	var wg sync.WaitGroup

	DetermCmd(usr, string(input), &wg)

	wg.Wait()
	repl(usr)

	recover()
}

// Instantiate objects and Start main loop. This is the function to
// start the application
func Run() {
	cli.Ready()
	defer cli.Restore()

	term := cli.NewTerminal()
	term.Clear()

	term.RawMode()
	defer term.Normal()

	usr := NewUser(term)

	logo(usr.InOut)

	repl(usr)
}

// Start Input Loop that after its done reading input and
// filling buffers, concurrently processes each.
func StartInputLoop(usr *User) line {
	if usr.InBody {
		printAndPrompt(&usr.InOut.term.Cursor, &usr.InOut.lines, usr.NestDepth)
	} else {
		printLineLogo(&usr.InOut.term.Cursor)
	}

	usr.InOut.term.Cursor.X = len(LINELOGO)

	for {
		if usr.InOut.InDebug {
			var dbwg sync.WaitGroup
			dbwg.Add(1)
			usr.InOut.Debugger.Ready <- &dbwg
			dbwg.Wait()
		}

		usr.InOut.done = EventChan(1)
		usr.InOut.read()

		var nlwg sync.WaitGroup
		nlwg.Add(1)

		var bufs = []buffer{&usr.InOut.Fbuf, &usr.InOut.Rbuf, &usr.InOut.Spbuf, &usr.InOut.Mvbuf, &usr.InOut.Wbuf}
		go ProccessBuffers(bufs, usr.InOut, &nlwg)

		nlwg.Wait()

		usr.Shader.FindLiterals()
		out := usr.Shader.Shade(string(usr.InOut.lines[usr.InOut.term.Cursor.Y]))

		RenderLine(&usr.InOut.term.Cursor, out)

		select {
		case <-usr.InOut.done:
			return usr.InOut.lines[usr.InOut.term.Cursor.Y-1]
		default:
			continue
		}

	}
}
