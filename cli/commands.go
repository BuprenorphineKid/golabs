package cli

import(
  "strings"
  "os"
  "os/exec"
  "fmt"
  "bufio"
)

func Eval() {
  proc := exec.Command("go", "run", ".labs/session/lab.go")
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

func DetermineCmd(lab *Lab, input string, usr *User) {
  switch {
  case strings.HasPrefix(input, "import"):
    go Import(lab, input)
  case strings.HasPrefix(input, ";eval"):
    Eval()
  case strings.HasPrefix(input, ";save"):
    Save()
  case strings.HasPrefix(input, ";help"):
    Help()
  default:
    go AddToMain(lab,input, usr)
  }
}

 func AddToMain(lab *Lab, line string, user *User) {
   l := *lab
   InsertString(l.Main, line + "\n", l.MainLine + user.CmdCount)

   user.addCmd(line)
 }
