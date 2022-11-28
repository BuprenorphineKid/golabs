package repl

import (
  "os"
  "bufio"
  "strings"
)

type User struct {
  CmdCount int
  Name string
  CmdHist []string
}

func NewUser() *User {
  h := make([]string, 1, 150)
  
  var u User
  u.CmdCount = 1
  u.CmdHist = h
  u.CmdHist[0] = "begin"

  return &u
}

func (u *User) setName(name string) {
  u.Name = name
}

func (u *User) addCmd(cmd string) {
  u.CmdHist = append(u.CmdHist, cmd)
  u.CmdCount++
}

func getInput(r *bufio.Reader) string {
  in, err := r.ReadString('\n')

  if err != nil {
    panic(err)
  }

  in = strings.TrimSpace(in)

  return in
}

func repl(r *bufio.Reader, lab *Lab, usr *User) {
  printLineLogo()
  input := getInput(r)

  DetermineCmd(lab, input, usr)

  repl(r, lab, usr)
}

func Run() {
  welcome()
  
  reader := bufio.NewReader(os.Stdin)
  lab := NewLab()
  usr := NewUser()
  
  repl(reader, lab, usr)
}

