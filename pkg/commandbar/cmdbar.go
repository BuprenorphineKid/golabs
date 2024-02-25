package commandbar

import (
	"fmt"
	"labs/pkg/cli"
	"labs/pkg/syntax"
)

type CommandBar struct {
	bar
	input string
}

var cursor = cli.NewCursor()

func NewCommandBar(h, w, x, y int, c, s string) *CommandBar {
	var cmd CommandBar
	cmd.bar = newBar(h, w, x, y, c, s)

	return &cmd
}

func (c CommandBar) Display() {
	c.bar.display()

	cursor.MoveTo(c.x+2, c.y+1)
	switch c.color {
	case "white":
		fmt.Print(syntax.OnWhite(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "black":
		fmt.Print(syntax.OnBlack(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "grey":
		fmt.Print(syntax.OnGrey(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "red":
		fmt.Print(syntax.OnRed(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "blue":
		fmt.Print(syntax.OnGrey(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "green":
		fmt.Print(syntax.OnGreen(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "yellow":
		fmt.Print(syntax.OnYellow(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "magenta":
		fmt.Print(syntax.OnMagenta(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	case "cyan":
		fmt.Print(syntax.OnCyan(syntax.Green(syntax.Italicized(" #Command# " + syntax.Magenta("<< ")))))
	}
}
