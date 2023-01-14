package repl

import (
	"fmt"
)

// Implimented to be s subject for observer pattern.
type subject interface {
	Register(Device)
	Remove(Device)
	Notify()
}

// Output Subject object for an observer pattern in which the Output obj will
// notify output devices to then be displayed.
type Output struct {
	line    string
	devices map[string]Device
	shader  Shader
}

// Global unexported instance of Output just to make it universally
// accessible. Seems to make the most sense to me
var output = newOutput(newHiLiter())

// Returns pointer to Output obj
func newOutput(s Shader) *Output {
	o := Output{}
	o.devices = make(map[string]Device)
	o.shader = s

	return &o
}

// Register method for aquiring new output devices.
func (o *Output) Register(n string, d Device) {
	o.devices[n] = d
}

// Remove method for deleting devices.
func (o *Output) Remove(d Device) {
	for k, v := range o.devices {
		if v == d {
			delete(o.devices, k)
		}
	}
}

// Notify all of your devices of the changes made to output.
func (o *Output) Notify() {
	for _, v := range o.devices {
		v.update(o.line)
	}
}

// Method of updating line to whats current, shading it, tben alerting
// the devices.
func (o *Output) SetLine(ln string) {
	o.line = o.shader.Shade(ln)
	o.Notify()
}

// Implemented to be eligible as an output device.
type Device interface {
	update(string)
}

// Displayable device
type Display interface {
	RenderLine(Cursor)
	PrintBuffer()
	PrintInPrompt(Cursor)
	PrintOutPrompt(Cursor)
	PrintAndPrompt(Cursor, *[]line, int)
}

// Default output device, could be whatever, like a log.
type Screen struct {
	prevLines []string
	Line      string
}

// Returns pointer to Screen obj
func newScreen() *Screen {
	s := Screen{}

	s.prevLines = make([]string, 0, 0)

	return &s
}

// Display method, implementation details dont matter that much as long as you
// make the incoming data your at some point.
func (s *Screen) update(line string) {
	s.Line = line
}

// Simply Print out the current string in the Line field.
func (s *Screen) PrintBuffer() {
	fmt.Print(s.Line)
}

// Moves cursor to beginning of line, cuts to the end and prints out
// current Line buffer, this is to update/redraw the lines as theyre being
// typed.
func (s *Screen) RenderLine(c Cursor) {
	c.Invisible()
	c.MoveTo(len(INPROMPT), c.GetY())
	c.CutRest()
	s.PrintBuffer()

	c.MoveTo(c.GetX(), c.GetY())
	c.Normal()
}

// Move Cursor to true beginning of current line, and prints the Colored
// INPROMPT (CINPROMPT) constant.
func (s *Screen) PrintInPrompt(c Cursor) {
	c.Invisible()

	c.MoveTo(0, c.GetY())
	fmt.Print(CINPROMPT)
	c.SetX(len(INPROMPT))

	c.Normal()
}

// Move Cursor to relative beginning of curent line, then prints thd Colored
// AND PROMPT (CANDPROMPT) constant.
func (s *Screen) PrintAndPrompt(c Cursor, ln *[]line, depth int) {
	c.Invisible()

	l := *ln

	c.MoveTo(len(INPROMPT)-len(ANDPROMPT), c.GetY())
	fmt.Print(CANDPROMPT)
	c.MoveTo(len(INPROMPT), c.GetY())
	c.SetX(len(INPROMPT))

	for j := 0; j < depth; j++ {
		l[c.GetY()] = l[c.GetY()].Tab(c.GetX())
		Tab()

		c.AddX(4)
		s.RenderLine(c)
	}

	c.Normal()
}

func (s *Screen) PrintOutPrompt(c Cursor) {
	c.Invisible()

	c.MoveTo(len(INPROMPT)-len(OUTPROMPT), c.GetY())
	fmt.Print(COUTPROMPT)
	c.MoveTo(len(INPROMPT), c.GetY())
	c.SetX(len(INPROMPT))

	c.Normal()
}
