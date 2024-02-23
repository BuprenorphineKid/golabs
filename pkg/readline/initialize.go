package readline

func Init() {
	out.Register("main", newScreen())
}
