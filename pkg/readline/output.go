package readline

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
	Devices map[string]Device
	shader  Shader
}

// Global unexported instance of Output just to make it universally
// accessible. Seems to make the most sense to me
var Out = newOutput(newHiLiter())

// Returns pointer to Output obj
func newOutput(s Shader) *Output {
	o := Output{}
	o.Devices = make(map[string]Device)
	o.shader = s

	return &o
}

// Register method for aquiring new output devices.
func (o *Output) Register(n string, d Device) {
	o.Devices[n] = d
}

// Remove method for deleting devices.
func (o *Output) Remove(d Device) {
	for k, v := range o.Devices {
		if v == d {
			delete(o.Devices, k)
		}
	}
}

// Notify all of your devices of the changes made to output.
func (o *Output) Notify() {
	for _, v := range o.Devices {
		v.Update(o.line)
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
	Update(string)
}

// Displayable device
type Display interface {
	RenderLine()
	PrintBuffer()
	PrintInPrompt()
	PrintOutPrompt()
	PrintAndPrompt(*[]line, int)
}

// Default output device, could be whatever, like a log.
type Screen struct {
	prevLines []string
	Line      string
}

// Returns pointer to Screen obj
func NewScreen() *Screen {
	s := Screen{}

	s.prevLines = make([]string, 0, 0)

	return &s
}

// Display method, implementation details dont matter that much as long as you
// make the incoming data your at some point.
func (s *Screen) Update(line string) {
	s.Line = line
}

// Simply Print out the current string in the Line field.
func (s *Screen) PrintBuffer() {
	fmt.Print(s.Line)
}

// Moves cursor to beginning of line, cuts to the end and prints out
// current Line buffer, this is to update/redraw the lines as theyre being
// typed.
func (s *Screen) RenderLine() {
	Term.Cursor.Invisible()
	Term.Cursor.MoveTo(len(INPROMPT), Term.Cursor.GetY())
	Term.Cursor.CutRest()
	s.PrintBuffer()

	Term.Cursor.MoveTo(Term.Cursor.GetX(), Term.Cursor.GetY())
	Term.Cursor.Normal()
}

// Move Cursor to true beginning of current line, and prints the Colored
// INPROMPT (CINPROMPT) constant.
func (s *Screen) PrintInPrompt() {
	Term.Cursor.Invisible()

	Term.Cursor.MoveTo(0, Term.Cursor.GetY())
	fmt.Print(CINPROMPT)
	Term.Cursor.SetX(len(INPROMPT))

	Term.Cursor.Normal()
}

// Move Cursor to relative beginning of curent line, then prints thd Colored
// AND PROMPT (CANDPROMPT) constant.
func (s *Screen) PrintAndPrompt(ln *[]line, depth int) {
	Term.Cursor.Invisible()

	l := *ln

	Term.Cursor.MoveTo(len(INPROMPT)-len(ANDPROMPT), Term.Cursor.GetY())
	fmt.Print(CANDPROMPT)
	Term.Cursor.MoveTo(len(INPROMPT), Term.Cursor.GetY())
	Term.Cursor.SetX(len(INPROMPT))

	for j := 0; j < depth; j++ {
		l[Term.Cursor.GetY()] = l[Term.Cursor.GetY()].Tab(Term.Cursor.GetX())
		Tab()

		Term.Cursor.AddX(4)
		s.RenderLine()
	}

	Term.Cursor.Normal()
}

func (s *Screen) PrintOutPrompt() {
	Term.Cursor.Invisible()

	Term.Cursor.MoveTo(len(INPROMPT)-len(OUTPROMPT), Term.Cursor.GetY())
	fmt.Print(COUTPROMPT)
	Term.Cursor.MoveTo(len(INPROMPT), Term.Cursor.GetY())
	Term.Cursor.SetX(len(INPROMPT))

	Term.Cursor.Normal()
}
