package repl

import (
	"fmt"
	"labs/syntax"
	"strings"
)

func DetermDecl(usr *User, inp string) {
	if usr.InBody == true {
		go Body(usr, inp)
		return
	}

	switch {
	case strings.HasPrefix(inp, "import"):
		go Import(usr.Lab, inp)
	case strings.HasPrefix(inp, "type"):
		go Type(usr, inp)
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
	if strings.Contains(bodyLine, "{") {
		usr.NestDepth += strings.Count(bodyLine, "{")
	}

	if strings.Contains(bodyLine, "}") {
		usr.NestDepth -= strings.Count(bodyLine, "}")
	}

	InsertString(usr.Lab.Main, bodyLine+"\n", usr.Lab.MainLine)
	usr.Lab.MainLine += 1

	if usr.NestDepth <= 0 {
		usr.InBody = false
	}
}
