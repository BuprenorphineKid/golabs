package labs

type History struct {
	count   int
	Entries []string
	Last    *string
}

// History constructor.
func NewHistory() *History {
	h := new(History)
	h.count = 1
	h.Entries = []string{}

	return h
}

// Increment CmdCount and add desired command to Hist.
func (h *History) Add(cmd string) {
	h.Entries = append(h.Entries, cmd)
	h.count++
	h.Last = &h.Entries[len(h.Entries)-1]
}
