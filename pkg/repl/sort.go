package repl

import (
	"regexp"
	"strings"
)

func RemovePrint(in string) string {
	if !strings.Contains(in, "print()") || !strings.Contains(in, "Print") {
		return in
	}

	re := regexp.MustCompile(`(fmt)?\.?[Pp]rint.*`)

	var repl string

	repl = re.ReplaceAllString(in, repl)

	return repl
}
