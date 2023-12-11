package readline

import (
	"regexp"
	"strings"

	"labs/pkg/syntax"
)

type Shader interface {
	//	FindLiterals()
	Shade(string) string
}

type HiLiter struct {
	List     []string
	done     []string
	strings  []string
	ints     []string
	keyWords map[string]string
	literals []string
}

func newHiLiter() *HiLiter {
	s := HiLiter{}

	s.done = make([]string, 0, 0)
	s.strings = make([]string, 0, 0)
	s.ints = make([]string, 0, 0)
	s.literals = make([]string, 0, 0)

	s.keyWords = map[string]string{
		"struct":     syntax.Blue("struct"),
		"interface":  syntax.Blue("interface"),
		"chan":       syntax.Yellow("chan"),
		"string":     syntax.Yellow("string"),
		"int":        syntax.Yellow("int"),
		"byte":       syntax.Yellow("byte"),
		"rune":       syntax.Yellow("rune"),
		"any":        syntax.Yellow("any"),
		"comparable": syntax.Yellow("comparable"),
		"complex128": syntax.Yellow("complex128"),
		"complex64":  syntax.Yellow("complex64"),
		"int8":       syntax.Yellow("int8"),
		"int16":      syntax.Yellow("int16"),
		"int32":      syntax.Yellow("int32"),
		"int64":      syntax.Yellow("int64"),
		"uint":       syntax.Yellow("uint"),
		"uint8":      syntax.Yellow("uint8"),
		"uint16":     syntax.Yellow("uint16"),
		"uint32":     syntax.Yellow("uint32"),
		"uint64":     syntax.Yellow("uint64"),
		"uintptr":    syntax.Yellow("uintptr"),
		"error":      syntax.Yellow("error"),
		"fost32":     syntax.Yellow("float32"),
		"float64":    syntax.Yellow("float64"),
		"bool":       syntax.Yellow("bool"),
		"nil":        syntax.Magenta("nil"),
		"for":        syntax.Red("for"),
		"if":         syntax.Red("if"),
		"switch":     syntax.Red("switch"),
		"else":       syntax.Red("else"),
		"select":     syntax.Red("select"),
		"recover":    syntax.Green("recover"),
		"panic":      syntax.Green("panic"),
		"make":       syntax.Green("make"),
		"copy":       syntax.Green("copy"),
		"new":        syntax.Green("new"),
		"append":     syntax.Green("append"),
		"len":        syntax.Green("len"),
		"complex":    syntax.Red("complex"),
		"imag":       syntax.Red("imag"),
		"cap":        syntax.Red("cap"),
		"delete":     syntax.Red("delete"),
		"print":      syntax.Red("print"),
		"println":    syntax.Red("println"),
		"real":       syntax.Red("real"),
		"close":      syntax.Red("close"),
		"go":         syntax.Cyan("go"),
		"return":     syntax.Red("return"),
		"range":      syntax.Red("range"),
		"func":       syntax.Red("func"),
		"type":       syntax.Red("type"),
		"true":       syntax.Magenta("true"),
		"false":      syntax.Magenta("false"),
		"iota":       syntax.Magenta("iota"),
		"import":     syntax.Cyan("import"),
		"package":    syntax.Cyan("package"),
		"var":        syntax.Blue("var"),
		"const":      syntax.Blue("const"),
	}

	s.List = []string{
		"struct",
		"interface",
		"chan",
		"string",
		"int",
		"byte",
		"rune",
		"any",
		"comparable",
		"complex128",
		"complex64",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"error",
		"float32",
		"float64",
		"bool",
		"nil",
		"for",
		"if",
		"switch",
		"else",
		"select",
		"recover",
		"panic",
		"make",
		"copy",
		"new",
		"append",
		"len",
		"complex",
		"imag",
		"cap",
		"delete",
		"print",
		"println",
		"real",
		"close",
		"go",
		"return",
		"range",
		"func",
		"type",
		"true",
		"false",
		"iota",
		"import",
		"package",
		"var",
		"const",
	}

	return &s
}

//func (s *HiLiter) FindLiterals() {
//	for _, v := range *s.lines {
//		s.strings = append(s.strings, syntax.Strings(string(v))...)
//		s.ints = append(s.ints, syntax.Ints(string(v))...)
//	}
//
//	s.literals = append(s.literals, s.strings...)
//	for i, v := range s.strings {
//		s.strings[i] = syntax.Green(v)
//	}
//
//	s.literals = append(s.literals, s.ints...)
//	for i, v := range s.ints {
//		s.ints[i] = syntax.Magenta(v)
//	}
//}

func (s *HiLiter) Shade(str string) string {
	var b = str

	for i, v := range s.List {
		re := regexp.MustCompile(`(^|\b)` + strings.TrimSpace(s.List[i]) + `(\b|$)`)
		b = re.ReplaceAllString(b, s.keyWords[v])
	}

	return b
}
