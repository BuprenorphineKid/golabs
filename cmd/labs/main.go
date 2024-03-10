package main

import (
	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"github.com/BuprenorphineKid/golabs/pkg/repl"
)

func main() {
	cli.Args()

	repl.Run()

}
