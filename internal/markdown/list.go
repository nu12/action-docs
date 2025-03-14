package markdown

type List struct {
	Items []string
}

func (l *List) Add(item string) *List {
	l.Items = append(l.Items, item)
	return l
}

func (l *List) String() string {
	var s = ""
	for _, item := range l.Items {
		s += "* " + item + "\n"
	}
	return s + "\n"
}
