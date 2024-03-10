/* This custom package which for me symbolizes my initiative in
 * starting to try to make Golang projects the RIGHT way.
 *
 * Scripts package is a simple API to access outside scripts
 * from inside your project. It has a few design innovations
 * by me but other than that, its just a wrapper for os/exec
 * for dead simple but consistent usage.
 *
 * Most of the functions will return a function closure of type
 * |func() error| so that the function can take in whatever
 * paramaters that youd like, and the Handler will still also
 * be able to call them on arrival.
 *
 *	+Basic Usage looks like as Follows+
 *
 * Exec >
 * -------
 * hndlr := scripts.NewHandler()
 * hndlr.Run()
 *
 * ruby := scripts.NewLanguage("ruby")
 * hndlr.Do <-scripts.Exec(ruby, "scripts/reset_server.rb")
 *
 * Eval >
 * -------
 *
 */
package scripts

import (
	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"log"
	"os/exec"
	"runtime/debug"
)

// The Handler type is a struct that.. well, handles your
// script executions, in such a way that makes it possible to
// do so concurrently.
type Handler struct {
	Do chan func() error
}

// Constructor for the Handler type. Call this for an pointer-ed
// instance.
func NewHandler() *Handler {
	s := new(Handler)
	s.Do = make(chan func() error)

	return s
}

// Run is what seperates "scripts" from "os/exec". You call
// this method on Handler to wait for functions on its
// Do channel, then after it recieves one, it immidiatly
// executes it no hesitation.
//
// This method is already nested in an anonymous goroutine
// as to not break the flow wherever youre calling it from. So
// theres no need to put go in front of calls to Run().
func (s *Handler) Run() {
	go func() {
		for {
			select {
			case f := <-s.Do:
				err := f()
				return
				if err != nil {
					cli.Restore()
					log.Fatalf(
						"\n\rfunc: |%s|\n\rerr: |%v|\n\r stack: |%s|\n\r",
						"Handler.Run() #"+f().Error(),
						err,
						debug.Stack(),
					)
				}
			}
		}
	}()
}

// Language couldve just easily been a custom string but it now
// acts as a representation of a scripting language that
// your desired script is wrote in.
//
// e.g. Lua, ruby, bash, etc.
type Language struct {
	name    string
	CallCmd exec.Cmd
}

// Languages constructor function.
func NewLanguage(lang string) Language {
	l := Language{name: lang}
	l.CallCmd = *exec.Command(lang)

	return l
}

// Exec() and Eval() are probably going to be the most used
// funcs in this pkg.
//
// Call this function, supplying a Language and a path to the
// script youd like to execute, and you will be returned a
// |func() error| that once called, will actually execute
// your script at the path you supplied.
func Exec(lang Language, script string) func() error {
	return func() error {
		lang.CallCmd.Args = append(lang.CallCmd.Args, script)

		err := lang.CallCmd.Run()
		if err != nil {
			return err
		}

		return nil
	}
}

// Exec() and Eval() are probably going to be the most used
// funcs in this pkg.
//
// Call this function, supplying a Language and a path to the
// script youd like to evaluate, and you will be returned a
// |chan []byte| where you can wait for the output of running
// your script and a |func() error| that once called,
// will actually execute your script at the path you
// supplied.
//
// The only difference between Exec() and Eval() are:
//
// With Eval() you get to see your output and with Exec(),
// you do not. Exec() is recommended if youre not keen
// on chanels and concurrency yet, as it will keep things
// a little simpler.
func Eval(lang Language, path string) (chan []byte, func() error) {
	output := make(chan []byte)

	return output, func() error {
		lang.CallCmd.Args = append(lang.CallCmd.Args, path)

		out, err := lang.CallCmd.Output()
		if err != nil {
			return err
		}

		output <- out

		return nil
	}
}
