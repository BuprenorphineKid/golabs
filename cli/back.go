//package main
package cli

import (
  //"fmt"
  "bufio"
  "io"
  "io/ioutil"
  "os"
  "strings"
)


type Lab struct {
  Main string
  Lines []string
  MainLine int
  ImportLine int
}

func NewLab() *Lab {
  l := Lab{}
  l.Main = project()
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

func project() string {
  content, _ := ioutil.ReadFile(template())
  os.Mkdir(".labs/session", 0777)
  proj, _ := os.Create(".labs/session/lab.go")
  defer proj.Close()
  proj.Write(content)

  return ".labs/session/lab.go"
}

func template() string {
  os.Mkdir(".labs", 0777)

  data := "\npackage main\n\nimport(\n\n)\n\nfunc main() {\n\n}\n"
  
  data = strings.TrimSpace(string(data))

  os.WriteFile(".labs/template", []byte(data), 0777)

  return ".labs/template"
}

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


///**
//* Insert string to n-th line of file.
//* If you want to insert a line, append newline '\n' to the end of the string.
//**/
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
