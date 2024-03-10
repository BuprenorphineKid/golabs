package eval

import (
	"fmt"
	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

type Report struct {
	Ok      bool
	Results string
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

	h, _ := os.UserHomeDir()
	e.file = h + "/.labs/session/eval.go"

	func() {

		e.Lock()
		cont, _ := os.ReadFile(path)
		_ = os.WriteFile(e.file, cont, 0777)
		e.Unlock()
	}()

	e.imports = exec.Command("goimports", "-w", e.file)
	e.imports.Stdin = os.Stdin

	e.format = exec.Command("gofmt", "-w", e.file)
	e.format.Stdin = os.Stdin

	e.run = exec.Command("go", "run", e.file)
	e.run.Stdin = os.Stdin

	return e
}

func (e *Evaluator) Exec(output chan Report) {
	e.Lock()
	defer e.Unlock()

	err := e.imports.Run()
	if err != nil {
		err := e.format.Run()
		if err != nil {
			if strings.Contains(err.Error(), "2") {
				output <- Report{Results: "", Ok: false}
				return
			}

			cli.NewTerminal().Normal()
			cli.Restore()

			log.Fatalf("\n\r|%s|\n\r%v", "Evaluator.Exec()", err)
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
			output <- Report{Results: "", Ok: false}
		}

		output <- Report{Results: fmt.Sprintf("Err: %v", strings.TrimSpace(string(f))), Ok: true}

		return
	}

	output <- Report{Results: string(res), Ok: true}

}
