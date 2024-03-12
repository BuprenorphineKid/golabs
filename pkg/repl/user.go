package repl

import (
	"sync"

	"github.com/BuprenorphineKid/golabs/pkg/labs"
	"github.com/BuprenorphineKid/golabs/pkg/readline"
)

// User Struct for keeping count of Hist CmdCount, yada yada.
// User holds all the main ingredients to run the show that
// are exvlusive to a user. Think env, files, input/output etc.
type User struct {
	Name     string
	Input    *readline.Input
	Lab      *labs.Lab
	done     chan struct{}
	FileLock sync.Locker
}

// Creates a new User object and returns a pointer to it.
func NewUser() *User {
	var u User

	u.Input = readline.NewInput()
	u.Lab = labs.NewLab()

	u.FileLock = new(sync.Mutex)

	return &u
}

// Set the Users Name.
func (u *User) setName(name string) {
	u.Name = name
}
