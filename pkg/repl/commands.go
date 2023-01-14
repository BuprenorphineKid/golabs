package repl

import (
	"bufio"
	"fmt"
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
		DetermDecl(usr, inp, m)
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

			c, err := os.ReadFile(".labs/session/lab.go")
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
	fmt.Println("Commands\n\r________\n\r';eval'  -  evaluate and print output of code so far\n\r';save'  -  save all of which you have just written to a new or existing File\n\r';help'  -  print this help message\n\r")
}
