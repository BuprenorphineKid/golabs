package repl

import(
  "strings"
  "os"
  "os/exec"
  "fmt"
  "bufio"
)

func DetermineCmd(lab *Lab, input string, usr *User) {
  switch {
  case strings.HasPrefix(input, "import"):
    go Import(lab, input)
  case strings.HasPrefix(input, ";eval"):
    Eval()
  case strings.HasPrefix(input, ";save"):
    Save()
  case strings.HasPrefix(input, ";help"):
    go Help()
  case strings.HasPrefix(input, "type"):
    go Type(lab, input)
  case strings.Contains(input, "struct") || 
      strings.HasPrefix(input, "type") && 
      strings.Contains(input, "struct"):
    go Struct(lab, input)
  default:
    go AddToMain(lab, input, usr)
  }
}

func Eval() {
  proc := exec.Command("go", "run", ".labs/session/lab.go")
  proc.Stderr = os.Stderr
  proc.Stdin = os.Stdin

  output, _ := proc.Output()

  fmt.Println(string(output))
}

func Import(lab *Lab, s string) {
  words := strings.Split(s, " ")
  pkg := strings.Trim(words[1], "\"")

  final := func () string {
    var p = []string{"\"", pkg, "\""}
    f := strings.Join(p, "")
    return f + "\n"
  }()

  InsertString(lab.Main, final, lab.ImportLine + 1)
  lab.MainLine++
}


func Type(lab *Lab, s string) {
  words := strings.Split(s, " ")
  newTypes := words[1:len(words) - 1]
  tp := words[len(words) - 1]

  for _, t := range newTypes {
    go func () {
      dec := fmt.Sprintf("type %s %s\n", t, tp)
      InsertString(lab.Main, dec, lab.MainLine)
      lab.MainLine++
    }()
  }
}

func Struct(lab *Lab, s string) {
  words := strings.Split(strings.TrimSpace(strings.Trim(s, "{}")), " ")
  var args []string

  if words[0] == "type" {
    args = words [3:]
  } else if words[0] == "struct" {
    args = words[2:]
  } 

  obj := words[1]

  f :=  make([]string, 0, len(args))
  t := make([]string, 0, len(args))
  
  for i, v := range args {  
    if i % 2 == 0 {
      f = append(f, v)
    } else {
      t = append(t, v)
    }
  }

  d := fmt.Sprintf("type %s struct {\n", obj)


  if len(f) > len(t) {
    f = f[:len(f) - 2]
  }
  if len(t) > len(f) { 
    t = t[:len(t) - 2]
  } 

  for k, _ := range f {  
    d += fmt.Sprintf("%s %s\n", f[k], t[k])
  }

  d += "}\n\n"

  InsertString(lab.Main, d, lab.MainLine - 1) 

  lab.MainLine += (3 + len(f))
}

func Save() {
  fmt.Print("Save File?\nno/Yes :")

  reader := bufio.NewReader(os.Stdin)
  resp := getInput(reader)

  switch {
  case resp == "y" || resp == "ye" || resp == "yes" || resp == "Y" || resp == "YE" || resp == "YES" || resp == "Ye" || resp == "Yes" || resp == "yES" || resp == "yeS":
    func () {
      fmt.Print("File Name :")

      name := getInput(reader)

      c, err := os.ReadFile(".labs/session/lab.go")
      if err != nil {
        panic(err)
      }
      
      file, _ := os.Create(name)
      file.Write(c)
      file.Close()
    }()
  default:
    return
  }
}
func Help() {
  fmt.Println("Commands\n________\n';eval'  -  evaluate and print output of code so far\n';save'  -  save all of which you have just written to a new or existing File\n';help'  -  print this help message")
}

func AddToMain(lab *Lab, line string, user *User) {
  l := *lab
  InsertString(l.Main, line + "\n", l.MainLine + user.CmdCount)

  user.addCmd(line)
}
