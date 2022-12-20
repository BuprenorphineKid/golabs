package repl

import (
	"fmt"
	"labs/cli"
	"strings"
)

const (
	LINELOGO  = "G-o-[-L-@-ß-$-] # "
	CLINELOGO = "\033[36;1;4;9mG\033[0;36;4;9m-\033[0;35;4;1mo\033[0;36;2;4;9m-\033[0;36;4;9m[\033[0;36;1;4;9m-\033[37;9;1mL\033[0;36;4;9m-\033[0;35;1;9m@\033[0;36;1;4;9m-\033[37;4;1mß\033[0;36;4;9m-\033[37;1;4m$\033[0;36;2;4;9m-\033[0;36;4;9;m]\033[0m \033[35;1;5m#\033[0m "

	LOGO = "\033[36;4;9m  __|  _ \\ |      \\   _ )  __|\033[0m\n\r\033[35;4;9m (_ | (   ||     _ \\  _ \\__ \\ \033[0m\n\r\033[36;1;4;9m\\___|\\___/____|_/  _\\___/____/\033[0m\n\r"

	ANDPROMPT  = "[-&?-] # "
	CANDPROMPT = "\033[36;1;4;9m[\033[36;4;2;9;3m-\033[0;35;4;9;1m&\033[35;9;1m?\033[0;36;1;4;9;3m-\033[0;36;4;9m]\033[0m \033[0;35;1m#\033[0m "
)

func printLineLogo(c *cli.Cursor) {
	fmt.Print(CLINELOGO)
	c.AddX(len(LINELOGO))
}

func printAndPrompt(c *cli.Cursor, l *[]line, depth int) {
	lines := *l

	c.MoveTo(len(LINELOGO)-len(ANDPROMPT), c.Y)
	fmt.Print(CANDPROMPT)

	for j := 0; j < depth; j++ {
		lines[c.Y] = lines[c.Y].Tab(c.X - len(LINELOGO))
		c.End(len(lines[c.Y]) + len(LINELOGO))
		c.AddX(8)
		RenderLine(c, string(lines[c.Y]))
	}
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

func PrintBuffer(c *cli.Cursor, out string) {
	fmt.Print(out)
}

func RenderLine(c *cli.Cursor, out string) {
	c.Invisible()
	c.MoveTo(len(LINELOGO), c.Y)
	c.CutRest()
	PrintBuffer(c, out)

	c.Normal()
}

func Backspace() {
	print("\b")
}
