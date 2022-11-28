package cli5

import (
  "flag"
  "os"
  "labs/repl"
)

type flag interface {}

var (
  last bool
  load string
)

func init() {
  flag.BoolVar(&last, "l", false, "start labs from previous session")
  flag.BoolVar(&last, "last", false, "start labs from previous session")
  flag.StringVar(&load, "L", "", "start labs with your own script")
  flag.StringVar(&load, "load", "", "start labs with your own script")
}

func Args() {
  flag.Parse
  
  if last == true {

  } else if load != "" {
    
  }
}
