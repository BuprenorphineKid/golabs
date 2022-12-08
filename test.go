package main

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"golang.org/x/term"
)

type Rbuf []byte

type ass interface{}

func test(a []ass) {
	for _, v := range a {
		fmt.Println(string(*v.(*Rbuf)))
		reflect.ValueOf(v).Elem().SetBytes(
			[]byte("dick"),
		)
		fmt.Println(string(*v.(*Rbuf)))
	}
}
func main() {
	state, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), state)

	s := newSome()

	var d = Rbuf("ass")
	var v = Rbuf("shit")

	var t = []ass{&d, &v}

	test(t)

	reflect.ValueOf(s).Elem().Field(2).SetBytes(
		[]byte("hello"),
	)
	s.loop()

}

type Some struct {
	writer io.Reader
	reader io.Writer
	Rbuf   Rbuf
}

func newSome() *Some {
	s := Some{
		reader: os.Stdin,
		writer: os.Stdout,
	}

	s.Rbuf = []byte("")

	return &s
}

func (s *Some) loop() {
	var buf [1]byte

	for {
		_, err := s.reader.(*os.File).Read(buf[:])

		if err != nil {
			panic(err)
		}

		if string(buf[0]) == "\x03" {
			break
		}

		s.Rbuf = append(s.Rbuf, buf[0])

		print(string(s.Rbuf))

		s.Rbuf = []byte("")
	}

	return
}
