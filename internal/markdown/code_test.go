package markdown

import (
	"testing"
)

func TestCode(t *testing.T) {
	c := Code("Hello")
	expected := "```\nHello\n```\n\n"
	result := c.String()
	if result != expected {
		t.Errorf("Code doesn't match. Got %q, want %q", result, expected)
	}
}
