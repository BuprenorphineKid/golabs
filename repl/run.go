package repl

import (
	"labs/cli"
	"os"
	"sync"
)

type Cursor interface {
	MoveTo(int, int)
	Home(int)
	End(int)
	CutRest()
	AddX(int)
	AddY(int)
	SetX(int)
	SetY(int)
	SavePos()
	RestorePos()
	Invisible()
	Normal()
	Up()
	Down()
	Left()
	Right()
}

// User Struct for keeping count of Hist CmdCount, yada yada.
// User holds all the main ingredients to run the show that
// are exvlusive to a user. Think env, files, input/output etc.
type User struct {
	CmdCount  int
	Name      string
	CmdHist   []string
	InBody    bool
	NestDepth int
	Input     *Input
	Lab       *Lab
	Logger    *Log
	done      chan struct{}
}

// Creates a new User object and returns a pointer to it.
func NewUser(t *cli.Terminal) *User {
	h := make([]string, 1, 150)

	var u User

	u.CmdCount = 1
	u.CmdHist = h
	u.CmdHist[0] = "begin"

	u.Input = NewInput(t)
	u.Lab = NewLab()

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

// Instantiate objects, Start the main loop, and start up the
// Evaluation Controller. This is the function you use to
// start the application
func Run() {
	cli.Ready()
	defer cli.Restore()

	term := cli.NewTerminal()
	term.Clear()
	term.RawMode()
	defer term.Normal()

	lgr := NewLog()

	controller := newController()
	usr := NewUser(term)
	usr.Logger = lgr

	output.Register("main", newScreen())

	logo(usr.Input)

	for len(controller.procs) == 0 {
		go repl(usr, controller)
	}

	controller.Run(usr.Input, &usr.Input.term.Cursor, usr.Logger)
}

// The main recursive loop at the top level.
func repl(usr *User, cntrl *Controller) {
	if usr.InBody {
		output.devices["main"].(Display).PrintAndPrompt(&usr.Input.term.Cursor, &usr.Input.lines, usr.NestDepth)
	} else {
		output.devices["main"].(Display).PrintInPrompt(&usr.Input.term.Cursor)
	}

	input := GetLine(usr.Input)

	DetermCmd(usr, string(input))

	runner := newgoRunner()
	var contents, _ = os.ReadFile(usr.Lab.Main)

	cntrl.Add(runner, contents)

	repl(usr, cntrl)
}

// Start Input Loop that after its done reading input and
// filling buffers, concurrently processes each.
func GetLine(i *Input) line {
	for {
		if i.InDebug {
			var dbwg sync.WaitGroup
			dbwg.Add(1)
			i.Debugger.Ready <- &dbwg
			dbwg.Wait()
		}

		i.done = EventChan(1)
		i.read()

		var nlwg sync.WaitGroup
		nlwg.Add(1)

		var bufs = []buffer{&i.Fbuf, &i.Rbuf, &i.Spbuf, &i.Mvbuf, &i.Wbuf}
		go ProccessBuffers(bufs, i, &i.term.Cursor, &nlwg)

		nlwg.Wait()

		output.SetLine(string(i.lines[i.term.Cursor.Y]))
		output.devices["main"].(Display).RenderLine(&i.term.Cursor)

		select {
		case <-i.done:
			return i.lines[i.term.Cursor.Y-1]
		default:
			continue
		}

	}
}
