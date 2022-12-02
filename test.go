package main

import (
  "fmt"
)

func main() {
  fmt.Printf("%#v %#v %#v %#v\n", "\x0a", "\x0b", "\x0c", "\x0a")
}
