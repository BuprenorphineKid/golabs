package repl

import (
  //"fmt"
  "bufio"
  "io"
  "io/ioutil"
  "os"
  "strings"
)


//
// Start the session with this one
//

func SetupSession(l *Lab) {
  temp := writeTemplate()

  writeSessionFile(l.LoadSessionContent(l.LoadContent))
}


//
// Main "Session" struct to hold state  
//

type Lab struct {
  Main string
  Lines []string
  MainLine int
  ImportLine int
}

type Contents struct {
  loaded []byte
  current []byte
}

//
// Lab Constructor
//

func NewLab() *Lab {
  l := Lab{}
  l.Main = l.p()
  l.Lines, _ =  file2lines(".labs/session/lab.go")

  var(
    ich = make(chan int)
    mch = make(chan int)
  )

  go func () {
    for i, s := range l.Lines {
      if s == "func main() {" {
        mch <- i
      }
    } 
  }()

  go func () {
    for i, s := range l.Lines {
      if s == "import(" {
        ich <- i
      }
    }
  }()

  l.MainLine = <- mch
  l.ImportLine = <- ich

  return &l
}


//
// Load it up
//

func (l *Lab) LoadSessioFile(file string) []byte {
  l.loadContent = ioutil.ReadFile(file)
  
  return l.loadContent
}


//
// Write session file, duh
//

func (l *Lab) writeSessionFile(content []byte) {
  os.Mkdir(".labs/session", 0777)
  
  proj, _ := os.Create(".labs/session/lab.go")
  defer proj.Close()
  
  proj.Write(content)
}


//
// Write the Template for session
// it may or may not be used dependin on cli flags
//

func writeTemplate() string {
  os.Mkdir(".labs", 0777)
 
  data := "\npackage main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n"

  os.WriteFile(".labs/template", []byte(data), 0777)

  return ".labs/template"
}


//
// Helper functions for inserting string/
// into a file
//

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


//
// Insert string to n-th line of file.
// If you want to insert a line, append newline '\n' to the end of the string.
//

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

  return ioutil.WriteFile(path, []byte(fileContent), 0644)
}
