package labs

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
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
	err := os.MkdirAll(".labs/session", 0777)
	if err != nil {
		panic(PERMERROR)
	}

	proj, _ := os.Create(".labs/session/lab.go")
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
	*History
}

// Lab Constructor
func NewLab() *Lab {
	l := Lab{}
	l.Main = ".labs/session/lab.go"
	l.Lines, _ = file2lines(".labs/session/lab.go")
	l.InBody = false
	l.Depth = 0
	l.History = NewHistory()

	var (
		ich  = make(chan int)
		mch  = make(chan int)
		done = make(chan struct{})
		wg   sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		for i, s := range l.Lines {
			if s == "func main() {" {
				mch <- i
				wg.Done()
				return
			}
		}
	}()

	go func() {
		for i, s := range l.Lines {
			if strings.HasPrefix(s, "import") {
				ich <- i
				wg.Done()
				return
			} else if i == len(l.Lines)-1 {
				wg.Done()
				done <- struct{}{}
			}
		}
	}()

loop:
	for {
		select {
		case l.MainLine = <-mch:
			continue
		case l.ImportLine = <-ich:
			continue
		case <-done:
			wg.Wait()
			break loop
		default:
		}

		if l.MainLine > 0 && l.ImportLine > 0 {
			break loop
		}
	}

	return &l
}

// Write the Template for session
// it may or may not be used dependin on cli flags
func writeTemplate() []byte {
	os.Mkdir(".labs", 0777)

	data := "package main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n"

	os.WriteFile(".labs/template", []byte(data), 0777)

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
