package cli

type Cursor struct {
	X int
	Y int
}

func newCursor() *Cursor {
	c := Cursor{}
	return &c
}

func (c *Cursor) SavePos() {
	print("\0337")
}

func (c *Cursor) RestorePos() {
	print("\0338")
}

func (c *Cursor) TrueHome() {
	print("\033[0;0H")
	c.X = 0
	c.Y = 0
}

func (c *Cursor) Home(prompt int) {
	print("\033[", c.Y, ";", prompt, "H")
}

func (c *Cursor) End(end int) {
	print("\033[", c.Y, ";", end, "H")
}

func (c *Cursor) Left() {
	print("\033[D")
}

func (c *Cursor) Right() {
	print("\033[C")
}

func (c *Cursor) Up() {
	print("\033[A")
}

func (c *Cursor) Down() {
	print("\033[B")
}

func (c *Cursor) GetX() int {
	return c.X
}

func (c *Cursor) GetY() int {
	return c.Y
}

func (c *Cursor) SetX(n int) {
	c.X = n
}

func (c *Cursor) AddX(n int) {
	c.X = c.X + n
}

func (c *Cursor) SetY(n int) {
	c.Y = n
}

func (c *Cursor) AddY(n int) {
	c.Y = c.Y + n
}

func (c *Cursor) MoveTo(x int, y int) {
	print("\033[", y, ";", x, "H")
}

func (c *Cursor) CutRest() {
	print("\033[0K")
}

func (c *Cursor) CutFirst() {
	print("\033[1K")
}

func (c *Cursor) CutLine() {
	print("\033[2K")
}

func (c *Cursor) Invisible() {
	print("\033[?25l")
}

func (c *Cursor) Normal() {
	print("\033[?12l\033[?25h")
}
