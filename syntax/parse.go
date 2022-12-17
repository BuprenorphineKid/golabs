package syntax

import (
	"regexp"
	"strings"
)

func FuncParts(s string) []string {

	recv := regexp.MustCompile(`^\(.+?\b.+?\)`)
	fnName := regexp.MustCompile(`^\b.+?\b|\B\s\b.+?\b`)
	param := regexp.MustCompile(`\(\)|\b\(.+?\b.+?\)`)
	retval := regexp.MustCompile(`\s\(.+?\)$|([[:alnum:]]|\*)*\b$`)

	rcv := strings.TrimSpace(recv.FindString(s))
	name := strings.TrimSpace(fnName.FindString(s))
	prm := strings.TrimSpace(param.FindString(s))
	ret := strings.TrimSpace(retval.FindString(s))

	return []string{rcv, name, prm, ret}
}

func TypeParts(s string) []string {
	name := regexp.MustCompile(`^\b.+?\b`)
	typ := regexp.MustCompile(`\s.+?$`)

	n := strings.TrimSpace(name.FindString(s))
	t := strings.TrimSpace(typ.FindString(s))

	return []string{n, t}
}
