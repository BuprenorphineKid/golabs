package repl

import (
  "fmt"
 // "bufio"
  //"os"
)
const(
LINELOGO = "(G-o-[-L-@-ß-$-]) # "
)

func printLineLogo() {
  fmt.Print(LINELOGO)
} 

func namePrompt() {
  fmt.Print("Enter UserName :  ")  
}

func welcome() {
  fmt.Println("Welcome to GoLabs, thank you for trying it out. This is a\nGolang based Repl so that prolly gives you a pretty good ides\n on how to use it and what not, but yeah just use it like you would any other\n Repl but with Golang code, and hopefully the rest should be self\n explanitory. type \";help\" for more info on commands and such.\nEnjoy.")
} 
