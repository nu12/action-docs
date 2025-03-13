package markdown

import (
	"testing"
)

func TestMarkdown(t *testing.T) {
	m := &Markdown{}
	h1 := H1("Hello")
	h2 := H2("World")
	m.Add(h1).Add(h2)
	expected := "# Hello\n\n## World\n\n"
	result := m.String()
	if result != expected {
		t.Errorf("Add doesn't match. Got %q, want %q", result, expected)
	}
}
