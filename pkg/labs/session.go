package labs

import (
	"bufio"
	"io"
	"os"
	"strings"
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
	print(h)
	err := os.MkdirAll(h+"/.labs/session", 0777)
	if err != nil {
		panic(PERMERROR)
	}

	proj, _ := os.Create(h + "/.labs/session/lab.go")
	defer proj.Close()

	proj.Write(content)
}

// Start the session with this one
func (c *Content) Setup() {
	c.writeSessionFile(c.Loaded)
}

// Main Session struct to hold state
type Lab struct {
	Main       string
	Lines      []string
	MainLine   int
	ImportLine int
	InBody     bool
	Depth      int
	History    *History
}

// Lab Constructor
func NewLab() *Lab {
	h, _ := os.UserHomeDir()
	l := Lab{}
	l.Main = h + "/.labs/session/lab.go"
	l.Lines, _ = file2lines(h + "/.labs/session/lab.go")
	l.InBody = false
	l.Depth = 0
	l.History = NewHistory()

	func() {
		for i, s := range l.Lines {
			if strings.HasPrefix(s, "func main()") {
				l.MainLine = i + 1
				return
			}
		}
	}()

	func() {
		for i, s := range l.Lines {
			if strings.HasPrefix(s, "import") {
				l.ImportLine = i
				return
			}
		}
	}()

	return &l

}

// Write the Template for session
// it may or may not be used dependin on cli flags
func writeTemplate() []byte {
	os.Mkdir("/.labs", 0777)

	h, _ := os.UserHomeDir()

	data := "package main\n\n\nimport(\n\n)\n\n\nfunc main() {\n\n\n}\n"

	os.WriteFile(h+"/.labs/template", []byte(data), 0777)

	return []byte(data)
}

// Helper functions for inserting string/
// into a file
func file2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	return linesFromReader(f)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// Insert string to n-th line of file.
// If you want to insert a line, append newline '\n' to the end of the string.
func InsertString(path, str string, index int) error {
	lines, err := file2lines(path)

	if err != nil {
		return err
	}

	fileContent := ""

	for i, line := range lines {
		if i == index {
			fileContent += str
		}

		fileContent += line
		fileContent += "\n"

	}

	return os.WriteFile(path, []byte(fileContent), 0644)
}

func Replace(path string, str string, index int) error {
	lines, err := file2lines(path)
	if err != nil {
		panic(err)
	}

	replacenent := strings.TrimSpace(str) + "\n"

	var fileContent string

	for k, v := range lines {
		if k == index {
			fileContent += replacenent
			continue
		}

		fileContent += v
		fileContent += "\n"
	}

	return os.WriteFile(path, []byte(fileContent), 0644)
}
