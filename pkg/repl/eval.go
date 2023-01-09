package repl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type printSlip struct {
	ok      bool
	results string
}

type Evaluator struct {
	imports *exec.Cmd
	format  *exec.Cmd
	run     *exec.Cmd
	file    string
}

func NewEvaluator(path string, m *sync.Mutex) *Evaluator {
	m.Lock()

	e := new(Evaluator)
	e.file = ".labs/session/eval.go"

	cont, _ := os.ReadFile(path)

	_ = os.WriteFile(e.file, cont, 0777)

	m.Unlock()

	e.imports = exec.Command("goimports", "-w", e.file)
	e.imports.Stdin = os.Stdin

	e.format = exec.Command("go", "fmt", "-w", e.file)
	e.format.Stdin = os.Stdin

	e.run = exec.Command("go", "run", e.file)
	e.run.Stdin = os.Stdin

	return e
}

func Eval(usr User, output chan printSlip, m *sync.Mutex) {
	eval := NewEvaluator(usr.Lab.Main, m)

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

		if strings.Contains(string(f), "not used") {
			output <- printSlip{results: "", ok: false}
		}

		output <- printSlip{results: fmt.Sprintf("Err: %v", strings.TrimSpace(string(f))), ok: true}
	}

	output <- printSlip{results: string(res), ok: true}
}
