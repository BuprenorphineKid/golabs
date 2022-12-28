package repl

import (
	"os"
	"os/exec"
)

type Evaluator interface {
	Eval([]byte, *Log) chan string
}

type goRunner struct {
	in   *os.File
	file string
}

func newgoRunner() *goRunner {
	g := goRunner{}
	g.in = os.Stdin
	g.file = ".labs/session/eval.go"

	return &g
}

func (g *goRunner) Eval(in []byte, l *Log) chan string {
	out := make(chan string)

	go func() {
		_ = os.WriteFile(g.file, in, 0777)

		goimp := exec.Command("goimports", "-w", g.file)
		goimp.Stdin = g.in
		goimp.Start()

		err := goimp.Wait()

		if err != nil {
			gofmt := exec.Command("gofmt", "-w", g.file)
			gofmt.Stdin = g.in
			gofmt.Start()

			err := gofmt.Wait()

			if err != nil {
				l.loggit(err)
			}
		}

		gorun := exec.Command("go", "run", g.file)
		gorun.Stdin = g.in

		output, err := gorun.Output()

		if err != nil {
			l.loggit(err)
		}

		out <- string(output)
	}()

	return out
}

type Controller struct {
	procs   map[Evaluator][]byte
	scripts []string
	results []string
}

func newController() *Controller {
	c := Controller{}
	c.procs = make(map[Evaluator][]byte)
	c.scripts = make([]string, 0, 0)
	c.results = make([]string, 0, 0)

	return &c
}

func (c *Controller) Add(e Evaluator, in []byte) {
	c.procs[e] = in
}

func (c *Controller) Remove(e Evaluator) {
	delete(c.procs, e)
}

func (c *Controller) Run(i *Input, curs Cursor, l *Log) {
	for {
		if len(c.procs) == 0 {
			continue
		}

		for k, v := range c.procs {
			go c.Catch(k.Eval(v, l), i, curs)
			c.scripts = append(c.scripts, string(v))
			c.Remove(k)
		}
	}
}

func (c *Controller) Catch(evalCh chan string, i *Input, curs Cursor) {
	select {
	case out := <-evalCh:
		c.results = append(c.results, out)
	}
}

func (c *Controller) ShowResults(curs Cursor) {
	if len(c.results) == 0 || c.results[len(c.results)-1] == "" {
		return
	}

	output.devices["main"].(Display).PrintInPrompt(curs)
	output.SetLine(c.results[len(c.results)-1])
	output.devices["main"].(Display).RenderLine(curs)
}
