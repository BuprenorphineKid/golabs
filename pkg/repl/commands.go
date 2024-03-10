package repl

import (
	"bufio"
	"fmt"
	"labs/pkg/scripts"
	"os"
)

func ExecuteCmd(usr *User, inp string) {
	switch inp {
	case "save", "write":
		SaveCmd()
		return
	case "help":
		HelpCmd()
		return
	case "reset", "clear", "new":
		ResetCmd()
		return
	default:
		return
	}

}

func SaveCmd() {
	fmt.Print("Save File?\nno/Yes :")

	var r = bufio.NewReader(os.Stdin)
	var buf []byte

	_, err := r.Read(buf)
	if err != nil {
		panic(err)
	}

	resp := string(buf)

	switch {
	case resp == "y" || resp == "ye" || resp == "yes" || resp == "Y" || resp == "YE" || resp == "YES" || resp == "Ye" || resp == "Yes" || resp == "yES" || resp == "yeS":
		func() {
			fmt.Print("File Name :")

			var buf []byte

			_, err := r.Read(buf)
			if err != nil {
				panic(err)
			}

			name := string(buf)

			h, _ := os.UserHomeDir()

			c, err := os.ReadFile(h + ".labs/session/lab.go")
			if err != nil {
				panic(err)
			}

			file, _ := os.Create(name)
			file.Write(c)
			file.Close()
		}()
	default:
		return
	}
}

func HelpCmd() {
	fmt.Println("Commands\n\r________\n\r';save'  -  save all of which you have just written to a new or existing File\n\r';help'  -  print this help message\n\r")
}

func ResetCmd() {
	scripter := scripts.NewHandler()
	scripter.Run()
	scripter.Do <- scripts.Exec(scripts.NewLanguage("bash"), "scripts/bash/new_session.sh")

	term.Clear()

	scrn.Buffer = make([]string, 0)

	usr.Input.CntrlCode <- 2
	usr = NewUser(term)

	InitializeUI()

}
