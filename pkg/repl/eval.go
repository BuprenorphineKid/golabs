package repl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

type Evaluator struct {
	imports *exec.Cmd
	format  *exec.Cmd
	run     *exec.Cmd
	file    string
}

func NewEvaluator(path string) *Evaluator {
	e := new(Evaluator)
	e.file = ".labs/session/eval.go"

	cont, _ := os.ReadFile(path)

	_ = os.WriteFile(e.file, cont, 0777)

	e.imports = exec.Command("goimports", "-w", path)
	e.imports.Stdin = os.Stdin

	e.format = exec.Command("go", "fmt", "-w", path)
	e.format.Stdin = os.Stdin

	e.run = exec.Command("go", "run", path)
	e.run.Stdin = os.Stdin

	return e
}

func Eval(usr User, output chan string) {
	eval := NewEvaluator(usr.Lab.Main)

	err := eval.imports.Run()

	if err != nil {
		err := eval.format.Run()

		if err != nil {
			log.Fatal(err)
		}

		eval.format.Wait()
	}

	eval.imports.Wait()

	res, err := eval.run.CombinedOutput()

	if err != nil {
		re := regexp.MustCompile(`(?s:.+:\d+:\s)`)

		var replacement = make([]byte, 0)
		f := re.ReplaceAll(res, replacement)

		output <- fmt.Sprintf("Err: %v", string(f))
	}

	output <- string(res)
}
