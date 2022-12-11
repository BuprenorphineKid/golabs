package repl

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	// "bufio"
)

func DetermineCmd(lab *Lab, inp line, usr *User, i *InOut) {
	switch {
	case strings.HasPrefix(string(inp), "import"):
		go Import(lab, string(inp))
	case strings.HasPrefix(string(inp), ";eval"):
		Eval()
	case strings.HasPrefix(string(inp), ";save"):
		Save(i)
	case strings.HasPrefix(string(inp), ";help"):
		go Help()
	case strings.HasPrefix(string(inp), "type"):
		go Type(lab, string(inp))
	case strings.Contains(string(inp), "struct") ||
		strings.HasPrefix(string(inp), "type") &&
			strings.Contains(string(inp), "struct"):
		go Struct(lab, string(inp))
	default:
		go AddToMain(lab, string(inp), usr)
	}
}

func Eval() {
	proc := exec.Command("go", "run", ".labs/session/lab.go")
	proc.Stderr = os.Stderr
	proc.Stdin = os.Stdin

	output, _ := proc.Output()

	fmt.Println(string(output))
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

func Save(i *InOut) {
	fmt.Print("Save File?\nno/Yes :")

	input := StartInputLoop(i)
	resp := *input

	switch {
	case resp == "y" || resp == "ye" || resp == "yes" || resp == "Y" || resp == "YE" || resp == "YES" || resp == "Ye" || resp == "Yes" || resp == "yES" || resp == "yeS":
		func() {
			fmt.Print("File Name :")

			input = StartInputLoop(i)
			name := string(*input)

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
	fmt.Println("Commands\n________\n';eval'  -  evaluate and print output of code so far\n';save'  -  save all of which you have just written to a new or existing File\n';help'  -  print this help message")
}

func AddToMain(lab *Lab, line string, user *User) {
	l := *lab
	InsertString(l.Main, line+"\n", l.MainLine+user.CmdCount)

	user.addCmd(line)
}
