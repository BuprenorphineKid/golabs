package labs

import (
	"fmt"
	"github.com/BuprenorphineKid/golabs/pkg/syntax"
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
		lab.Add(inp)
	}

}

func (lab *Lab) Add(inp string) {
	InsertString(lab.Main, inp+"\n", lab.MainLine+lab.History.count)

	lab.History.Add(inp)
}

func (lab *Lab) Replace(inp string, pos int) {
	Replace(lab.Main, inp+"\n", lab.MainLine+pos)

}

func Import(lab *Lab, s string) {
	words := strings.Split(s, " ")
	pkg := strings.Trim(words[1], "\"")

	final := func() string {
		var p = []string{"\"", pkg, "\""}
		f := strings.Join(p, "")
		return f + "\n"
	}()

	InsertString(lab.Main, final, lab.ImportLine)
	lab.MainLine++
}

func Type(lab *Lab, s string) {
	trimmed := strings.TrimSpace(strings.Trim(s, "{}type"))

	parts := syntax.TypeParts(trimmed)

	if parts[1] == "struct" || parts[1] == "interface" {
		dec := fmt.Sprintf("type %s %s {\n", parts[0], parts[1])
		InsertString(lab.Main, dec, lab.MainLine-1)

		lab.InBody = true
		lab.Depth++
		lab.MainLine++
	} else {
		dec := fmt.Sprintf("type %s %s{}\n", parts[0], parts[1])
		InsertString(lab.Main, dec, lab.MainLine-1)
		lab.MainLine++
	}
}

func Func(lab *Lab, s string) {
	trimmed := strings.TrimSpace(strings.Trim(s, "{}func"))

	parts := syntax.FuncParts(trimmed)

	decl := fmt.Sprintf("func %s %s%s %s {\n", parts[0], parts[1], parts[2], parts[3])

	InsertString(lab.Main, decl, lab.MainLine-1)
	lab.MainLine++
	lab.InBody = true
	lab.Depth++
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

	InsertString(lab.Main, space+strings.TrimLeft(bodyLine+"\n", " "), lab.MainLine-1)
	lab.MainLine++

	if lab.Depth <= 0 {
		lab.InBody = false
	}
}
