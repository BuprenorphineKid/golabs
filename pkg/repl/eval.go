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
	sync.Locker
	imports *exec.Cmd
	format  *exec.Cmd
	run     *exec.Cmd
	file    string
}

func NewEvaluator(path string) *Evaluator {
	e := new(Evaluator)
	e.Locker = new(sync.Mutex)
	e.file = ".labs/session/eval.go"

	func() {

		e.Lock()
		cont, _ := os.ReadFile(path)
		_ = os.WriteFile(e.file, cont, 0777)
		e.Unlock()
	}()

	e.imports = exec.Command("goimports", "-w", e.file)
	e.imports.Stdin = os.Stdin

	e.format = exec.Command("go", "fmt", "-w", e.file)
	e.format.Stdin = os.Stdin

	e.run = exec.Command("go", "run", e.file)
	e.run.Stdin = os.Stdin

	return e
}

func (e *Evaluator) Exec(output chan printSlip) {
	e.Lock()
	defer e.Unlock()

	err := e.imports.Run()

	if err != nil {

		err := e.format.Run()

		if err != nil {
			log.Fatal(err)
		}

		e.format.Wait()

	}

	e.imports.Wait()

	res, err := e.run.CombinedOutput()

	if err != nil {
		re := regexp.MustCompile(`(?s:.+:\d+:\s)`)

		var replacement = make([]byte, 0)
		f := re.ReplaceAll(res, replacement)

		if strings.Contains(string(f), "not used") {
			output <- printSlip{results: "", ok: false}
		}

		output <- printSlip{results: fmt.Sprintf("Err: %v", strings.TrimSpace(string(f))), ok: true}

		return
	}

	output <- printSlip{results: string(res), ok: true}

}
