package markdown

import (
	"testing"
)

func TestTable(t *testing.T) {
	header := Header{"Name", "Age"}
	rows := []Row{
		{"Alice", "20"},
		{"Bob", "30"},
	}
	table := Table{header, rows}

	expected := "|Name|Age|\n|---|---|\n|Alice|20|\n|Bob|30|\n\n"

	result := table.String()
	if result != expected {
		t.Errorf("Table doesn't match. Got %q, want %q", result, expected)
	}
}

func TestSort(t *testing.T) {
	header := Header{"Name", "Age"}
	rows := []Row{
		{"Bob", "30"},
		{"Alice", "20"},
	}
	table := Table{header, rows}

	expected := "|Name|Age|\n|---|---|\n|Alice|20|\n|Bob|30|\n\n"

	result := table.Sort(0).String()
	if result != expected {
		t.Errorf("Table doesn't match. Got %q, want %q", result, expected)
	}
}
