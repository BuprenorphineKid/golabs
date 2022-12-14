package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	// "bufio"
)

func DetermCmd(usr *User, inp string) {
	var wg sync.WaitGroup

	switch inp {
	case ";eval":
		wg.Add(1)
		go usr.Eval.Evaluate(usr.InOut, &wg)
		wg.Wait()
	case ";save":
		Save()
	case ";help":
		go Help()
	default:
		DetermDecl(usr, inp)
	}
}

type Eval struct {
	stdin   io.Reader
	stderr  io.Writer
	file    string
	LastOut []byte
}

func NewEval() *Eval {
	e := Eval{
		stdin:  os.Stdin,
		stderr: os.Stderr,
		file:   ".labs/session/lab.go",
	}

	return &e
}

func (e *Eval) Evaluate(i *InOut, wg *sync.WaitGroup) {
	var err error

	imp := exec.Command("goimports", "-w", e.file)
	imp.Stderr = e.stderr
	imp.Stdin = e.stdin

	err = imp.Start()
	if err != nil {
		fmt := exec.Command("go", "fmt", e.file)
		fmt.Stderr = e.stderr
		fmt.Stdin = e.stdin

		err = fmt.Start()
		if err != nil {
			panic(err)
		}

		fmt.Wait()
	}

	proc := exec.Command("go", "run", e.file)
	proc.Stderr = e.stderr
	proc.Stdin = e.stdin

	imp.Wait()

	e.LastOut, err = proc.Output()

	if err != nil {
		panic("Unfortunately Go is not installed, please install to run your code")
	}

	fmt.Println("\r" + string(e.LastOut) + "\r")

	func() {
		l := len(e.LastOut) / i.term.Cols

		parts := strings.Split(string(e.LastOut), "\n")

		l += len(parts)

		i.AddLines(l + 1)
		i.term.Cursor.AddY(l)
	}()

	wg.Done()
}

func Save() {
	fmt.Print("Save File?\nno/Yes :")

	var r = bufio.NewReader(os.Stdin)
	var buf []byte

	_, err := r.Read(buf)
	if err != nil {
		panic(err)
	}

	resp := string(buf)

	switch {
	case resp == "y" || resp == "ye" || resp == "yes" || resp == "Y" || resp == "YE" || resp == "YES" || resp == "Ye" || resp == "Yes" || resp == "yES" || resp == "yeS":
		func() {
			fmt.Print("File Name :")

			var buf []byte

			_, err := r.Read(buf)
			if err != nil {
				panic(err)
			}

			name := string(buf)

			c, err := os.ReadFile(".labs/session/lab.go")
			if err != nil {
				panic(err)
			}

			file, _ := os.Create(name)
			file.Write(c)
			file.Close()
		}()
	default:
		return
	}
}

func Help() {
	fmt.Println("Commands\n\r________\n\r';eval'  -  evaluate and print output of code so far\n\r';save'  -  save all of which you have just written to a new or existing File\n\r';help'  -  print this help message\n\r")
}
