package repl

import (
	"labs/cli"
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

// The main recursive loop at the top level.
func repl(usr *User) {
	input := StartInputLoop(usr.InOut)

	DetermCmd(usr, string(input))

	repl(usr)
}

// Instantiate objects and Start main loop. This is the function to
// start the application
func Run() {
	cli.Ready()

	term := cli.NewTerminal()
	term.Clear()
	term.RawMode()

	usr := NewUser(term)

	logo(usr.InOut)

	repl(usr)
}
