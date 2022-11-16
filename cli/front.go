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

  return in + "\n"
}

func logo() {
  fmt.Print("(G-o-[-L-@-ß-$-]) # ")
}

func repl(r *bufio.Reader, lab *Lab, count *int) {
  l := *lab
  c := *count
  
  logo()
  input := getInput(r)
  InsertString(l.Main, input, l.MainLine + c)

  c++
  
  repl(r, lab, &c)
}

func Run() {
  reader := bufio.NewReader(os.Stdin)

  lab := NewLab()

  var count int = 1
  
  repl(reader, lab, &count)
  
}
