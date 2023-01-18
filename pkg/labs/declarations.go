package labs

import (
	"fmt"
	"labs/pkg/syntax"
	"strings"
	"sync"
)

func DetermDecl(lab *Lab, inp string, m sync.Locker) {
	m.Lock()
	defer m.Unlock()

	if lab.InBody == true {
		Body(lab, inp)
		return
	}

	switch {
	case strings.HasPrefix(inp, "import"):
		Import(lab, inp)
	case strings.HasPrefix(inp, "type"):
		Type(lab, inp)
	case strings.HasPrefix(inp, "func"):
		Func(lab, inp)
	default:
		AddToMain(lab, inp)
	}

}

func AddToMain(lab *Lab, inp string) {
	InsertString(lab.Main, inp+"\n", lab.MainLine+lab.count)

	lab.AddCmd(inp)
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
	trimmed := strings.TrimSpace(strings.Trim(s, "{}type"))

	parts := syntax.TypeParts(trimmed)

	if parts[1] == "struct" || parts[1] == "interface" {
		dec := fmt.Sprintf("type %s %s {\n", parts[0], parts[1])
		InsertString(lab.Main, dec, lab.MainLine)

		lab.InBody = true
		lab.Depth += 1
		lab.MainLine += 1
	} else {
		dec := fmt.Sprintf("type %s %s\n", parts[0], parts[1])
		InsertString(lab.Main, dec, lab.MainLine)
		lab.MainLine += 1
	}

}

func Func(lab *Lab, s string) {
	trimmed := strings.TrimSpace(strings.Trim(s, "{}func"))

	parts := syntax.FuncParts(trimmed)

	decl := fmt.Sprintf("func %s %s%s %s {\n", parts[0], parts[1], parts[2], parts[3])

	InsertString(lab.Main, decl, lab.MainLine)
	lab.MainLine += 1
	lab.InBody = true
	lab.Depth += 1
}

func Body(lab *Lab, bodyLine string) {
	if !lab.InBody {
		return
	}

	if strings.Contains(bodyLine, "{") {
		lab.Depth += strings.Count(bodyLine, "{")
	}

	if strings.Contains(bodyLine, "}") {
		lab.Depth -= strings.Count(bodyLine, "}")
	}

	var space string
	for i := 0; i <= lab.Depth; i++ {
		space += "    "
	}

	InsertString(lab.Main, space+strings.TrimLeft(bodyLine+"\n", " "), lab.MainLine)
	lab.MainLine += 1

	if lab.Depth <= 0 {
		lab.InBody = false
	}
}
