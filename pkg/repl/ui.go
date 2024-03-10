package repl

import (
	"labs/pkg/readline"
)

func InitializeUI() {
	readline.Logo(usr.Input)
	readline.Init()

	scrn.Win.LoadScreen()
	scrn.Win.Fill()
	scrn.Win.Draw()

}
