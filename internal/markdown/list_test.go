package markdown

import (
	"testing"
)

func TestList(t *testing.T) {
	list := &List{}
	list.Add("item1").Add("item2")

	if len(list.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(list.Items))
	}

	if list.Items[0] != "item1" || list.Items[1] != "item2" {
		t.Errorf("items not added correctly, got %v", list.Items)
	}

	expected := "* item1\n* item2\n\n"
	result := list.String()

	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}
