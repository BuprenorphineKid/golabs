package repl

import (
	"fmt"
	"strings"
)

const (
	LINELOGO = "G-o-[-L-@-ß-$-] # "

	CLINELOGO = "\033[31mG\033[0;30;1m-\033[0;31mo\033[0;30;1m-\033[0;30;1m[\033[0;30;1m-\033[31;1mL\033[0;30;1m-\033[0;32m@\033[0;30;1m-\033[31;1mß\033[0;30;1m-\033[31;1m$\033[0;30;1m-\033[0;30;1m]\033[0m \033[36;1m#\033[0m "

	LOGO = "\033[36;1;4;9m  __|  _ \\ |      \\   _ )  __|\033[0m\n\r\033[36;1;4;9m (_ | (   ||     _ \\  _ \\__ \\ \033[0m\n\r\033[36;1;4;9m\\___|\\___/____|_/  _\\___/____/\033[0m\n\r"
)

func printLineLogo(i *InOut) {
	fmt.Print(CLINELOGO)
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

func PrintBuffer(i *InOut) {
	fmt.Print(
		string([]byte(i.lines[i.term.Cursor.Y])[i.term.Cursor.X-len(LINELOGO):]),
	)
}

func Refresh(i *InOut) {
	i.term.Cursor.CutRest()
	PrintBuffer(i)

	i.term.Cursor.MoveTo(
		i.term.Cursor.X+1,
		i.term.Cursor.Y,
	)
}

func Backspace() {
	print("\b")
}
