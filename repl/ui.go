package repl

import (
	"fmt"
	"strings"
)

const (
	LINELOGO = "(G-o-[-L-@-ß-$-]) # "

	LOGO = "  __|  _ \\ |      \\   _ )  __|\n\r (_ | (   ||     _ \\  _ \\__ \\\n\r\\___|\\___/____|_/  _\\___/____/\n\r"
)

func printLineLogo(i *InOut) {
	fmt.Print(LINELOGO)
	i.term.Cursor.AddX(len(LINELOGO))
}

func namePrompt(i *InOut) {
	fmt.Print("Enter UserName :  ")
	i.term.Cursor.AddX(len("Enter UserName :  "))
}

func logo(i *InOut) {
	fmt.Println(LOGO)

	parts := strings.Split(LOGO, "\n")

	i.term.Cursor.AddY(len(parts) + 1)
	i.AddLines(len(parts) + 1)
}
