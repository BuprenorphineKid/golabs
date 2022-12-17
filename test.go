package main

import (
	"fmt"
	"labs/syntax"
)

type myString string

type tr struct {
	mob int
	yes bool
}

func newtr() *tr {
	t := tr{}
	t.mob = 178
	t.yes = false

	return &t
}

func main() {
	a := 25

	b := "twenty-five"

	j := newtr()

	var c myString = "25"

	z := []byte("controller")

	x := syntax.Black(syntax.OnGreen(*j))

	d := syntax.OnBlue(syntax.Magenta(a))

	f := syntax.Blue(b)

	e := syntax.Cyan(syntax.OnBlack(c))

	s := syntax.Red(syntax.OnCyan(string(z)))

	fmt.Printf("%s\n%s\n%s\n%s\n%s\n", d, f, e, s, x)
}
