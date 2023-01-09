package repl

import (
	"fmt"
	"labs/pkg/syntax"
	"strings"
	"sync"
)

func DetermDecl(usr *User, inp string, m *sync.Mutex) {
	m.Lock()

	if usr.InBody == true {
		Body(usr, inp)
		return
	}

	switch {
	case strings.HasPrefix(inp, "import"):
		Import(usr.Lab, inp)
	case strings.HasPrefix(inp, "type"):
		Type(usr, inp)
	case strings.HasPrefix(inp, "func"):
		Func(usr, inp)
	default:
		AddToMain(usr, inp)
	}

	m.Unlock()
}

func AddToMain(usr *User, inp string) {

	InsertString(usr.Lab.Main, inp+"\n", usr.Lab.MainLine+usr.CmdCount)

	usr.addCmd(inp)
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

func Type(usr *User, s string) {
	trimmed := strings.TrimSpace(strings.Trim(s, "{}type"))

	parts := syntax.TypeParts(trimmed)

	if parts[1] == "struct" || parts[1] == "interface" {
		dec := fmt.Sprintf("type %s %s {\n", parts[0], parts[1])
		InsertString(usr.Lab.Main, dec, usr.Lab.MainLine)

		usr.InBody = true
		usr.NestDepth += 1
		usr.Lab.MainLine += 1
	} else {
		dec := fmt.Sprintf("type %s %s\n", parts[0], parts[1])
		InsertString(usr.Lab.Main, dec, usr.Lab.MainLine)
		usr.Lab.MainLine += 1
	}
}

func Func(usr *User, s string) {
	trimmed := strings.TrimSpace(strings.Trim(s, "{}func"))

	parts := syntax.FuncParts(trimmed)

	decl := fmt.Sprintf("func %s %s%s %s {\n", parts[0], parts[1], parts[2], parts[3])

	InsertString(usr.Lab.Main, decl, usr.Lab.MainLine)
	usr.Lab.MainLine += 1
	usr.InBody = true
	usr.NestDepth += 1
}

func Body(usr *User, bodyLine string) {
	if !usr.InBody {
		return
	}

	if strings.Contains(bodyLine, "{") {
		usr.NestDepth += strings.Count(bodyLine, "{")
	}

	if strings.Contains(bodyLine, "}") {
		usr.NestDepth -= strings.Count(bodyLine, "}")
	}

	var space string
	for i := 0; i <= usr.NestDepth; i++ {
		space += "    "
	}

	InsertString(usr.Lab.Main, space+strings.TrimLeft(bodyLine+"\n", " "), usr.Lab.MainLine)
	usr.Lab.MainLine += 1

	if usr.NestDepth <= 0 {
		usr.InBody = false
	}
}
