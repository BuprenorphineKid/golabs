package window

import (
	"fmt"
	"labs/pkg/cli"
	"labs/pkg/syntax"
	"strings"
)

type Screen struct {
	win     *Window
	curs    *cli.Cursor
	results []string
}

func NewScreen(w *Window, c *cli.Cursor) *Screen {
	s := Screen{}
	s.win = w
	s.curs = c
	s.results = make([]string, 0, 0)

	return &s
}

func (s *Screen) Wrap(buf string) {
	if strings.Contains(buf, "\n") {

		out := strings.Split(buf, "\n")

		for _, v := range out {
			w := s.win.Width

			if len(v) > w {
				s.results = append(s.results, syntax.OnGrey(syntax.White(v[:w])))

				s.results = append(s.results, syntax.OnGrey(syntax.White(v[w:])))

				continue
			}

			s.results = append(s.results, syntax.OnGrey(syntax.White(v)))
		}
	} else {
		s.results = append(s.results, syntax.OnGrey(syntax.White(buf)))
	}
}

func (s *Screen) Display() {
	s.curs.SavePos()

	s.win.Fill()
	defer s.win.Draw()

	ycurs := s.win.Y + 1
	xcurs := s.win.X + 3

	for _, v := range s.results {
		s.curs.MoveTo(xcurs, ycurs)

		fmt.Println(v)
		ycurs++
	}

	s.results = nil
	s.results = make([]string, 0)

	s.curs.RestorePos()
}
