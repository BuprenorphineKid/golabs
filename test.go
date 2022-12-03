package main

import "fmt"

func main() {
	fmt.Printf("home: %x\n\nend: %x\n", "\033[H", "\033[F")
	s := "jnn c c"
	s -= "c"
}
