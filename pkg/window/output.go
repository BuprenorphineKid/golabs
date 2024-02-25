package window

import (
	"fmt"
	"labs/pkg/cli"
	"labs/pkg/syntax"
	"slices"
	"strings"
)

type Screen struct {
	win    *Window
	curs   *cli.Cursor
	Buffer []string
}

func NewScreen(w *Window, c *cli.Cursor) *Screen {
	s := Screen{}
	s.win = w
	s.curs = c
	s.Buffer = make([]string, 0, 0)

	return &s
}

func (s *Screen) Reset() {
	s.Buffer = make([]string, 0)
}

func (s *Screen) Scroll() {
	s.Buffer = s.Buffer[1:]
	slices.Clip(s.Buffer)
}

func (s *Screen) Wrap(buf string) {
	if strings.Contains(buf, "\n") {

		out := strings.Split(buf, "\n")

		for _, v := range out {
			w := s.win.Width

			if len(v) > w {
				s.Buffer = append(s.Buffer, syntax.OnGrey(syntax.White(v[:w])))

				s.Buffer = append(s.Buffer, syntax.OnGrey(syntax.White(v[w:])))

				continue
			}

			s.Buffer = append(s.Buffer, syntax.OnGrey(syntax.White(v)))
		}
	} else {
		s.Buffer = append(s.Buffer, syntax.OnGrey(syntax.White(buf)))
	}

}

func (s *Screen) TrimSpace() {
	if len(s.Buffer) > 0 {
		return
	}
	for _, v := range s.Buffer {

		if strings.Contains(v, "\n") {
			slices.Delete(s.Buffer, len(s.Buffer), len(s.Buffer))
			slices.Clip(s.Buffer)
		}
	}
}

func (s *Screen) Display() {
	s.win.Fill()

	ycurs := s.win.Y + 1
	xcurs := s.win.X + 3

	for _, v := range s.Buffer {
		s.curs.MoveTo(xcurs, ycurs)

		fmt.Println(v)
		ycurs++
	}
	s.win.Draw()

	s.Reset()
}
