package labs

type History struct {
	count   int
	entries []string
	Last    *string
}

// History constructor.
func NewHistory() *History {
	h := new(History)
	h.count = 1
	h.entries = []string{}

	return h
}

// Increment CmdCount and add desired command to Hist.
func (h *History) AddCmd(cmd string) {
	h.entries = append(h.entries, cmd)
	h.count++
	h.Last = &h.entries[len(h.entries)-1]
}
