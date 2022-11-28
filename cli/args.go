package cli

import (
  "flag"
  "labs/repl"
)

var (
  last bool
  load string
)

func init() {
  flag.BoolVar(&last, "l", false, "[last] start labs from previous session")
  flag.StringVar(&load, "L", "", "[Load] start labs with your own script")
}

func Args(c *repl.Content) {  
  flag.Parse()

  if last == true {
    c.Load(".labs/session/lab.go")
  } else if load != "" {
    c.Load(load)
  } else if load != "" && last == true {
    panic("Contradicting flag options both on")
  }
}
