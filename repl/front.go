package repl

import (
	"labs/cli"
)

// User Struct for keeping count of Hist CmdCount, yada yada.
type User struct {
	CmdCount int
	Name     string
	CmdHist  []string
}

// Creates a new User object and returns a pointer to it.
func NewUser() *User {
	h := make([]string, 1, 150)

	var u User
	u.CmdCount = 1
	u.CmdHist = h
	u.CmdHist[0] = "begin"

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
func repl(i *InOut, lab *Lab, usr *User) {
	input := StartInputLoop(i)

	DetermineCmd(lab, *input, usr, i)

	repl(i, lab, usr)
}

// Instantiate objects and Start main loop.
func Run() {
	cli.Ready()

	term := cli.NewTerminal()
	term.Clear()
	term.RawMode()

	inout := newInOut(term)

	welcome(inout)

	lab := NewLab()
	usr := NewUser()

	repl(inout, lab, usr)
}
