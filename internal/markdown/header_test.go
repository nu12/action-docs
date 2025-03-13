package markdown

import (
	"testing"
)

func TestH1(t *testing.T) {
	h := H1("Hello")
	expected := "# Hello\n\n"
	result := h.String()
	if result != expected {
		t.Errorf("H1 doesn't match. Got %q, want %q", result, expected)
	}
}

func TestH2(t *testing.T) {
	h := H2("Hello")
	expected := "## Hello\n\n"
	result := h.String()
	if result != expected {
		t.Errorf("H2 doesn't match. Got %q, want %q", result, expected)
	}
}
