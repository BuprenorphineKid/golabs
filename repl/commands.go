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

func DetermineCmd(inp string, usr *User) {
	switch {
	case strings.HasPrefix(inp, "import"):
		go Import(usr.Lab, inp)
	case strings.HasPrefix(inp, ";eval"):
		var wg sync.WaitGroup

		wg.Add(1)
		usr.Eval.Evaluate(&wg)

		wg.Wait()
	case strings.HasPrefix(inp, ";save"):
		Save()
	case strings.HasPrefix(inp, ";help"):
		go Help()
	case strings.HasPrefix(inp, "type"):
		go Type(usr.Lab, inp)
	case strings.Contains(inp, "struct"):
		go Struct(usr.Lab, inp)
	default:
		go AddToMain(usr, inp)
	}
}

type Eval struct {
	stdin   io.Reader
	stderr  io.Writer
	Process *exec.Cmd
	LastOut []byte
}

func NewEval() *Eval {
	e := Eval{
		stdin:   os.Stdin,
		stderr:  os.Stderr,
		Process: exec.Command("go", "run", ".labs/session/lab.go"),
	}

	e.Process.Stderr = e.stderr
	e.Process.Stdin = e.stdin

	return &e
}

func (e *Eval) Evaluate(wg *sync.WaitGroup) {
	var err error

	e.LastOut, err = e.Process.Output()

	if err != nil {
		panic(err)
	}

	fmt.Println("\r" + string(e.LastOut))

	wg.Done()
}

func Import(lab *Lab, s string) {
	words := strings.Split(s, " ")
	pkg := strings.Trim(words[1], "\"")

	final := func() string {
		var p = []string{"\"", pkg, "\""}
		f := strings.Join(p, "")
		return f + "\n"
	}()

	InsertString(lab.Main, final, lab.ImportLine+1)
	lab.MainLine++
}

func Type(lab *Lab, s string) {
	words := strings.Split(s, " ")
	newTypes := words[1 : len(words)-1]
	tp := words[len(words)-1]

	go func() {
		for _, t := range newTypes {
			dec := fmt.Sprintf("type %s %s\n", t, tp)
			InsertString(lab.Main, dec, lab.MainLine)
			lab.MainLine++
		}
	}()
}

func Struct(lab *Lab, s string) {
	words := strings.Split(strings.TrimSpace(strings.Trim(s, "{}")), " ")
	var args []string

	if words[0] == "type" {
		args = words[3:]
	} else if words[0] == "struct" {
		args = words[2:]
	}

	obj := words[1]

	f := make([]string, 0, len(args))
	t := make([]string, 0, len(args))

	for i, v := range args {
		if i%2 == 0 {
			f = append(f, v)
		} else {
			t = append(t, v)
		}
	}

	d := fmt.Sprintf("type %s struct {\n", obj)

	if len(f) > len(t) {
		f = f[:len(f)-2]
	}
	if len(t) > len(f) {
		t = t[:len(t)-2]
	}

	for k, _ := range f {
		d += fmt.Sprintf("%s %s\n", f[k], t[k])
	}

	d += "}\n\n"

	InsertString(lab.Main, d, lab.MainLine-1)

	lab.MainLine += (3 + len(f))
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

func AddToMain(usr *User, row string) {
	InsertString(usr.Lab.Main, row+"\n\r", usr.Lab.MainLine+usr.CmdCount)

	usr.addCmd(row)
}
