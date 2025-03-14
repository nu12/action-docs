package markdown

type Element interface {
	String() string
}

type Markdown struct {
	Elements []Element
}

func (m *Markdown) Add(e Element) *Markdown {
	m.Elements = append(m.Elements, e)
	return m
}

func (m *Markdown) String() string {
	s := ""
	for _, e := range m.Elements {
		s += e.String()
	}
	return s
}
