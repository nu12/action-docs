package markdown

import "fmt"

type Hyperlink struct {
	URL  string
	Text string
}

func (h *Hyperlink) String() string {
	return fmt.Sprintf("[%s](%s)", h.Text, h.URL)
}
