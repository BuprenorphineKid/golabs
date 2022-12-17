package syntax

import "fmt"

type Blank interface {
	~string | ~int | any
}

type Color int

const CSI string = "\033["

const (
	FG int = 3
	BG int = 4
)

const (
	BLACK Color = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

const (
	m          = "m"
	END string = "\033[0m"
)

func Black[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, BLACK, m, word, END)
}

func Red[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, RED, m, word, END)
}

func Green[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, GREEN, m, word, END)
}

func Yellow[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, YELLOW, m, word, END)
}

func Blue[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, BLUE, m, word, END)
}

func Magenta[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, MAGENTA, m, word, END)
}

func Cyan[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, CYAN, m, word, END)
}

func White[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, FG, WHITE, m, word, END)
}

func OnBlack[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, BLACK, m, word, END)
}

func OnRed[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, RED, m, word, END)
}

func OnGreen[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, GREEN, m, word, END)
}

func OnYellow[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, YELLOW, m, word, END)
}

func OnBlue[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, BLUE, m, word, END)
}

func OnMagenta[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, MAGENTA, m, word, END)
}

func OnCyan[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, CYAN, m, word, END)
}

func OnWhite[T Blank](word T) string {
	return fmt.Sprintf("%s%d%d%s%v%s", CSI, BG, WHITE, m, word, END)
}
