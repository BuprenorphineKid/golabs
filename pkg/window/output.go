package window

import (
	"fmt"
	"github.com/BuprenorphineKid/golabs/pkg/cli"
	"github.com/BuprenorphineKid/golabs/pkg/syntax"
	"slices"
	"strings"
)

type Screen struct {
	Win     *Window
	curs    *cli.Cursor
	Buffer  []string
	history []string
}

func NewScreen(w *Window, c *cli.Cursor) *Screen {
	s := Screen{}
	s.Win = w
	s.curs = c
	s.Buffer = make([]string, 0, 0)

	return &s
}

func (s *Screen) Reset() {
	s.Buffer = make([]string, 0)
}

func (s *Screen) Scroll() {

	xs := len(s.Buffer) - (s.Win.Height - 2)

	s.history = s.Buffer[:xs]
	s.Buffer = s.Buffer[xs:]

	slices.Clip(s.Buffer)
}

func (s *Screen) Wrap(buf string) {
	if strings.Contains(buf, "\n") {

		out := strings.Split(buf, "\n")

		for _, v := range out {
			w := s.Win.Width

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
	if len(s.Buffer) == 0 {
		return
	}

	if s.Buffer[len(s.Buffer)] == "\n" {
		slices.Delete(s.Buffer, len(s.Buffer), len(s.Buffer))
		slices.Clip(s.Buffer)
	}
}

func (s *Screen) Display() {
	s.Win.Fill()

	ycurs := s.Win.Y + 1
	xcurs := s.Win.X + 3

	for _, v := range s.Buffer {
		s.curs.MoveTo(xcurs, ycurs)

		fmt.Println(v)
		ycurs++
	}
	s.Win.Draw()

	s.Reset()
}
