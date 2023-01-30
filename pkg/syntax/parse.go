package syntax

import (
	"log"
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

	for _, v := range parts {
		_, err = strconv.Atoi(v)
		if err == nil {
			i = append(i, v)
		}
	}

	return i
}

func IsFuncCall(s string) bool {
	if !strings.Contains(s, "(") && !strings.Contains(s, ")") {
		return false
	}

	match, err := regexp.MatchString(
		`\s*\b([A-Za-z]+[A-Za-z0-9_]*?\b)?\.?[A-Za-z]+[A-Za-z0-9_]*\b\(.*\)`,
		s,
	)
	if err != nil {
		log.Fatalf(
			"syntax.IsFuncCall()\n\rErr: %v\n\r",
			err,
		)
	}

	if !match {
		return false
	}

	types := []string{
		"string",
		"[]string",
		"bool",
		"[]byte",
		"byte",
		"rune",
		"[]rune",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uintptr",
		"struct{}",
	}

	for i := range types {
		if strings.HasPrefix(strings.TrimSpace(s), types[i]) {
			return false
		}
	}

	return true
}
