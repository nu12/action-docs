package markdown

import "sort"

type Table struct {
	Header Header
	Rows   []Row
}

type Header []string
type Row []string

func (t *Table) AddRow(r Row) {
	t.Rows = append(t.Rows, r)
}

func (t *Table) String() string {
	var table string

	table += "|"
	for _, h := range t.Header {
		table += h + "|"
	}
	table += "\n"

	table += "|"
	for range t.Header {
		table += "---|"
	}
	table += "\n"

	for _, r := range t.Rows {
		table += "|"
		for _, c := range r {
			table += c + "|"
		}
		table += "\n"
	}

	return table + "\n"
}

func (t *Table) Sort(columnPosition int) *Table {
	sort.Slice(t.Rows, func(i, j int) bool {
		return t.Rows[i][columnPosition] < t.Rows[j][columnPosition]
	})
	return t
}
