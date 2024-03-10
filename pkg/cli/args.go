package cli

import (
	"flag"
	"os"
)

var (
	last bool
	load string
)

// content struct
type Content struct {
	Loaded []byte
}

func NewContent() *Content {
	var c Content
	c.Loaded = writeTemplate()

	return &c
}

// Load it up
func (c *Content) Load(file string) {
	l, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	c.Loaded = l
}

// Write session file, duh
func (c *Content) writeSessionFile(content []byte) {
	h, _ := os.UserHomeDir()
	os.Mkdir(h+"/.labs/session", 0777)

	proj, _ := os.Create(h + "/.labs/session/lab.go")
	defer proj.Close()

	proj.Write(content)
}

// Start the session with this one
func (c *Content) Setup() {
	c.writeSessionFile(c.Loaded)
}

// Write the Template for session
// it may or may not be used dependin on cli flags
func writeTemplate() []byte {
	h, _ := os.UserHomeDir()

	os.Mkdir(h+"/.labs", 0777)

	data := "package main\n\n\nimport(\n\n)\n\n\nfunc main() {\n\n\n}\n"

	os.WriteFile(h+"/.labs/template", []byte(data), 0777)
	return []byte(data)
}

func init() {
	flag.BoolVar(&last, "l", false, "[last] start labs from previous session")
	flag.StringVar(&load, "L", "", "[Load] start labs with your own script")
}

func Args() {
	flag.Parse()

	h, _ := os.UserHomeDir()

	c := NewContent()

	if last == true {
		c.Load(h + "/.labs/session/lab.go")
	} else if load != "" {
		c.Load(load)
	} else if load != "" && last == true {
		panic("Contradicting flag options both on")
	}

	c.Setup()
}
