package repl

import (
	"fmt"
	"strings"
)

const (
	INPROMPT  = "G-o-[-L-@-ß-$-] <> "
	CINPROMPT = "\033[36;1;4;9mG\033[0;36;4;9m-\033[0;35;4;1mo\033[0;36;2;4;9m-\033[0;36;4;9m[\033[0;36;1;4;9m-\033[37;9;1mL\033[0;36;4;9m-\033[0;35;1;9m@\033[0;36;1;4;9m-\033[37;4;1mß\033[0;36;4;9m-\033[37;1;4m$\033[0;36;2;4;9m-\033[0;36;4;9;m]\033[0m \033[35;1;9;5m<>\033[0m "

	LOGO = "\033[36;4;9m  __|  _ \\ |      \\   _ )  __|\033[0m\n\r\033[35;4;9m (_ | (   ||     _ \\  _ \\__ \\ \033[0m\n\r\033[36;1;4;9m\\___|\\___/____|_/  _\\___/____/\033[0m\n\r"

	ANDPROMPT  = "[-&?-] << "
	CANDPROMPT = "\033[37;1;4;9m[\033[36;4;9;3m-\033[0;36;1;4;m&\033[0;35;4;1m?\033[0;30;4;9;3m-\033[0;30;4;9;2m]\033[0m \033[0;35;9m>>\033[0m "

	OUTPROMPT  = ">>"
	COUTPROMPT = "\033[35;9;1m<<\033[0m "
)

func logo(i *Input) {
	fmt.Println(LOGO)

	parts := strings.Split(LOGO, "\n")

	i.term.Cursor.AddY(len(parts) + 1)
	i.AddLines(len(parts) + 1)
}

func Backspace() {
	print("\b")
}

func Tab() {
	print("    ")
}
