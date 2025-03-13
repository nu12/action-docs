package markdown

import (
	"testing"
)

func TestP(t *testing.T) {
	p := P("Hello")
	expected := "Hello\n\n"
	result := p.String()
	if result != expected {
		t.Errorf("P doesn't match. Got %q, want %q", result, expected)
	}
}
