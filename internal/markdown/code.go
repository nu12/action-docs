package markdown

type Code string

func (c Code) String() string {
	return "```\n" + string(c) + "\n```\n\n"
}
