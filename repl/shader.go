package repl

type Shader struct {
	lines    *[]line
	done     []string
	strings  []string
	ints     []string
	types    []string
	ctrl     []string
	funcs    []string
	keyWords []string
}

func newShader(l *[]line) *Shader {
	s := Shader{
		lines: l,
	}

	s.todo = make([]string, 0, 0)
	s.done = make([]string, 0, 0)
	s.strings = make([]string, 0, 0)
	s.ints = make([]string, 0, 0)

	s.types = []string{
		syntax.Blue("struct"),
		syntax.Blue("interface"),
		syntax.Yellow("chan"),
		syntax.Yellow("strings"),
		syntax.Yellow("int"),
		syntax.Yellow("byte"),
		syntax.Yellow("rune"),
		syntax.Yellow("any"),
		syntax.Yellow("comparable"),
		syntax.Yellow("complex128"),
		syntax.Yellow("complex64"),
		syntax.Yellow("int8"),
		syntax.Yellow("int16"),
		syntax.Yellow("int32"),
		syntax.Yellow("int64"),
		syntax.Yellow("uint"),
		syntax.Yellow("uint8"),
		syntax.Yellow("uint16"),
		syntax.Yellow("uint32"),
		syntax.Yellow("uint64"),
		syntax.Yellow("uintptr"),
		syntax.Yellow("error"),
		syntax.Yellow("float32"),
		syntax.Yellow("float64"),
		syntax.Yellow("bool"),
		syntax.Magenta("nil"),
	}

	s.ctrl = []string{
		syntax.Red("for"),
		syntax.Red("if"),
		syntax.Red("switch"),
		syntax.Red("else"),
		syntax.Red("select"),
	}

	s.funcs = []string{
		syntax.Green("recover"),
		syntax.Green("panic"),
		syntax.Green("make"),
		syntax.Green("copy"),
		syntax.Green("new"),
		syntax.Green("append"),
		syntax.Green("len"),
		syntax.Red("complex"),
		syntax.Red("imag"),
		syntax.Red("cap"),
		syntax.Red("delete"),
		syntax.Red("print"),
		syntax.Red("println"),
		syntax.Red("real"),
		syntax.Red("close"),
		syntax.Cyan("go"),
	}

	s.keyWords = []string{
		syntax.Red("return"),
		syntax.Red("range"),
		syntax.Red("func"),
		syntax.Red("type"),
		syntax.Magenta("true"),
		syntax.Magenta("false"),
		syntax.Magenťa("iota"),
		syntax.Cyan("import"),
		syntax.Cyan("package"),
	}

	return &l
}

func (s *Shader) Parse(l line) {
	for _, v := range s.lines {
		s.strings = append(s.strings, syntax.Strings(v)...)
		s.ints = append(s.ints, syntax.Ints(v)...)
	}

	for i, v := range &s.strings {
		s.strings[i] = syntax.Green(v)
	}

	for i, v := range &s.ints {
		s.ints[i] = syntax.Magrnta(v)
	}
}

func (s *Shader) Inject(buf) {

}
