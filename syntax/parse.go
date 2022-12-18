package syntax

import (
	"regexp"
	"strconv"
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

func Strings(s string) []string {
	count := strings.Count(s, "\"") / 2

	if count%2 == 1 {
		idx := strings.LastIndex(s, "\"")
		s = string([]byte(s)[:idx])
		count--
	}

	sli := make([]string, count)
	str := regexp.MustCompile(`\B\".*?\"\B`)

	var st string

	for i := 0; i < count; i++ {
		st = strings.TrimSpace(str.FindString(s))

		sli = append(sli, st)

		s = strings.Replace(s, st, "", 1)
	}

	return sli
}

func Ints(s string) []string {
	parts := strings.Split(s, "")

	var i []string
	var err error

	for _, v := range b {
		_, err = strconv.Atoi(v)
		if err == nil {
			i = append(i, v)
		}
	}

	return i
}
