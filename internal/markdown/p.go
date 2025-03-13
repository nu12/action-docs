package markdown

type P string

func (p P) String() string {
	return string(p) + "\n\n"
}
