package repl

import (
	"log"
	"os"
)

type logger interface {
	Println(...any)
	Fatal(...any)
}

type Log struct {
	file *os.File
	path string
	name string
	logger
}

func NewLog() *Log {

	l := new(Log)
	l.path = ".labs/log/"
	_ = os.MkdirAll(l.path, os.ModeDir|0777)

	l.name = "session.txt"
	l.file, _ = os.OpenFile(l.path+l.name, os.O_CREATE|os.O_APPEND|os.O_SYNC, 0777)
	l.logger = log.New(
		l.file,
		"Oops![-L-@-B-$-] : ",
		log.Llongfile|log.Ldate|log.Ltime,
	)

	return l
}

func (l *Log) loggit(a ...any) {
	l.logger.Println(a...)
}

func (l *Log) Stoppit(a ...any) {
	l.logger.Fatal(a...)
}
