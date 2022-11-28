package main

import(
  "labs/repl"
  "labs/cli"
)

func main() {
  c := repl.NewContent()
  cli.Args(c)
  c.Setup()

  repl.Run()

}
