package repl

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

func RemovePrint(in string) string {
	if !strings.Contains(in, "print()") || !strings.Contains(in, "Print") {
		return in
	}

	re := regexp.MustCompile(`(fmt)?\.?[Pp]rint.*`)

	var repl string

	repl = re.ReplaceAllString(in, repl)

	return repl
}

func inferVar(str string, bufs [2][]string) {
	first, last, ok := strings.Cut(str, ":=")

	if ok == true {
		f := regexp.MustCompile(`\b[a-zA-z_]+?[a-zA-Z0-9_]*?\b\s$`)
		l := regexp.MustCompile(`^\s.\b+\(.*\)\B|^\s.+?\b`)

		bufs[0] = append(bufs[0], f.FindString(first))
		bufs[1] = append(bufs[1], l.FindString(last))
	}

	if strings.Contains(last, ":=") {
		inferVar(last, bufs)
	}

	return
}

func ParseVars(s string) {
	if !strings.ContainsAny(s, "var:=") {
		return
	}

	kk := make([]string, 0, 0)
	vv := make([]string, 0, 0)

	var bufs = [2][]string{kk, vv}

	go inferVar(s, bufs)

	for i := range bufs[0] {
		vars = append(vars, NewVar(bufs[0][i], bufs[1][i]))
	}
}

type Var struct {
	Name  string
	Value interface{}
	Type  string
}

var vars = make([]*Var, 0, 0)

func NewVar(name string, val ...interface{}) *Var {
	v := Var{Name: name, Value: val}

	return &v
}

func (v *Var) Typed() string {
	if v.Type != "" {
		return v.Type
	}

	v.Type = fmt.Sprintf("%T", v.Value)

	return v.Type
}

func (v *Var) Simplify() {
	m := new(sync.Mutex)

	for i, v := range vars {
		file, _ := os.Create(".labs/session/vars/")
		file.Write([]byte(v.Value.(string)))

		e := NewEvaluator(".labs/session/vars/", m)
		e
	}
}
