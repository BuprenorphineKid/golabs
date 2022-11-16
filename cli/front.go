package cli 

import (
  "fmt"
  "os"
  "bufio"
  "strings"
)
/*
type User struct {
  
}
*/
func getInput(r *bufio.Reader) string {
  in, err := r.ReadString('\n')

  if err != nil {
    panic(err)
  }

  in = strings.TrimSpace(in)

  return in
}

func logo() {
  fmt.Print("(L@ß$) % ")
}

func repl(r *bufio.Reader, lab *Lab) {
  l := *lab
  
  logo()
  input := getInput(r)
  InsertString(l.Main, input, l.MainLine + 1)
  
  repl(r, lab)
}

func Run() {
  reader := bufio.NewReader(os.Stdin)

  lab := NewLab()
  
  repl(reader, lab)
  
}
