package main

import (
	"labs/pkg/cli"
	"labs/pkg/repl"
)

func main() {
	cli.Args()

	repl.Run()

}
