package workflow

import (
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/markdown"
)

func TestWorkflows(t *testing.T) {
	const errorf = "Error: %v. \nExpected: %v \nGot: %v\nFrom: %v"
	tests := []struct {
		name     string
		given    []Workflow
		expected string
	}{
		{
			name: "Placeholder",
			given: []Workflow{
				{
					Name:        "A",
					Description: "Test workflows",
					Filename:    ".github/workflows/a.yml",
				},
			},
			expected: "0ffb9812157409a3206fb7c8ffd85f92",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := Workflows{
				Workflows: []Workflow{},
				Content:   markdown.List{},
			}
			for _, w := range tt.given {
				ws.AddWorkflow(&w)
			}
			md := ws.String()
			got := helper.Hash(md)
			if tt.expected != got {
				t.Errorf(errorf, "mismatch", tt.expected, got, md)
			}
		})
	}
}
