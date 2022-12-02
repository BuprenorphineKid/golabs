package main

import (
	"labs/cli"
	"labs/repl"
)

func main() {
	cli.Args()

	repl.Run()

}
