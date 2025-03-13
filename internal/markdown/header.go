package markdown

type H1 string
type H2 string
type H3 string

func (h H1) String() string {
	return "# " + string(h) + "\n\n"
}
func (h H2) String() string {
	return "## " + string(h) + "\n\n"
}
func (h H3) String() string {
	return "### " + string(h) + "\n\n"
}
