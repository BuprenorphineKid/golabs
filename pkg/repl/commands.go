package repl

import (
	"bufio"
	"fmt"
	"labs/pkg/labs"
	"os"
	"sync"
)

func DetermCmd(usr *User, inp string, m sync.Locker) {
	switch inp {
	case "!save":
		Save()
	case "!help":
		go Help()
	default:
		labs.DetermDecl(usr.Lab, inp, m)
	}

}

func ExecuteCmd(usr *User, inp string) {
	switch inp {
	case "!save":
		Save()
	case "!help":
		go Help()
	default:
	}

}

func Save() {
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

func Help() {
	fmt.Println("Commands\n\r________\n\r';save'  -  save all of which you have just written to a new or existing File\n\r';help'  -  print this help message\n\r")
}
