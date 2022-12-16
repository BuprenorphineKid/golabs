package repl

import (
	"fmt"
	"regexp"
	"strings"
)

func DetermDecl(usr *User, inp string) {
	if usr.InFunc == true {
		go Body(usr, inp)
		return
	}

	switch {
	case strings.HasPrefix(inp, "import"):
		go Import(usr.Lab, inp)
	case strings.HasPrefix(inp, "type"):
		go Type(usr.Lab, inp)
	case strings.Contains(inp, "struct"):
		go Struct(usr.Lab, inp)
	case strings.HasPrefix(inp, "func"):
		go Func(usr, inp)
	default:
		go AddToMain(usr, inp)
	}
}

func AddToMain(usr *User, row string) {
	InsertString(usr.Lab.Main, row+"\n", usr.Lab.MainLine+usr.CmdCount)

	usr.addCmd(row)
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

func Func(usr *User, s string) {
	trimmed := strings.TrimSpace(strings.Trim(s, "{}func"))

	recv := regexp.MustCompile(`^\(.+?\b.+?\)`)
	fnname := regexp.MustCompile(`^\b.+?\b|\B\s\b.+?\b`)
	param := regexp.MustCompile(`\(\)|\b\(.+?\b.+?\)`)
	retval := regexp.MustCompile(`\s\(.+?\)$|([[:alnum:]]|\*)*\b$`)

	rcv := strings.TrimSpace(recv.FindString(trimmed))
	name := strings.TrimSpace(fnName.FindString(trimmed))
	prm := strings.TrimSpace(param.FindString(trimmed))
	ret := strings.TrimSpace(retval.FindString(trimmed))

	decl := fmt.Sprintf("func %s %s%s %s {", rcv, name, prm, ret)

	InsertString(usr.Lab.Main, decl, usr.Lab.MainLine-1)
	usr.Lab.MainLine += 1
	usr.InFunc = true
	usr.NestDepth += 1
}

func Body(usr *User, bodyLine string) {
	if strings.Contains(bodyLine, "{") {
		usr.NestDepth += strings.Count(bodyLine, "{")
	}

	if strings.Contains(bodyLine, "}") {
		usr.NestDepth -= strings.Count(bodyLine, "}")
	}

	InsertString(usr.Lab.Main, bodyLine+"\n", usr.Lab.MainLine-1)
	usr.Lab.MainLine += 1

	if usr.NestDepth <= 0 {
		usr.InFunc = false
	}
}
