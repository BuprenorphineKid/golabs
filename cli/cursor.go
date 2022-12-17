package cli

type cursor struct {
	X int
	Y int
}

func newCursor() *cursor {
	c := cursor{}
	return &c
}

func (c *cursor) SavePos() {
	print("\0337")
}

func (c *cursor) RestorePos() {
	print("\0338")
}

func (c *cursor) TrueHome() {
	print("\033[0;0H")
	c.X = 0
	c.Y = 0
}

func (c *cursor) Home(prompt int) {
	print("\033[", c.Y, ";", prompt, "H")
}

func (c *cursor) End(end int) {
	print("\033[", c.Y, ";", end, "H")
}

func (c *cursor) Left() {
	print("\033[D")
}

func (c *cursor) Right() {
	print("\033[C")
}

func (c *cursor) Up() {
	print("\033[A")
}

func (c *cursor) Down() {
	print("\033[B")
}

func (c *cursor) AddX(n int) {
	c.X = c.X + n
}

func (c *cursor) AddY(n int) {
	c.Y = c.Y + n
}

func (c *cursor) MoveTo(x int, y int) {
	print("\033[", y, ";", x, "H")
}

func (c *cursor) CutRest() {
	print("\033[0K")
}

func (c *cursor) CutFirst() {
	print("\033[1K")
}

func (c *cursor) CutLine() {
	print("\033[2K")
}

func (c *cursor) Invisible() {
	print("\033[?25l")
}

func (c *cursor) Normal() {
	print("\033[?12l\033[?25h")
}
