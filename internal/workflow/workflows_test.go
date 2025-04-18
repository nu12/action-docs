package workflow

import (
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/markdown"
)

func TestWorkflows(t *testing.T) {
	tests := []struct {
		name         string
		given        []Workflow
		expectedHash string
	}{
		{
			name: "One workflow",
			given: []Workflow{
				{
					Name:        "Workflow A",
					Description: "Test workflows A",
					Filename:    ".github/workflows/a.yml",
				},
			},
			expectedHash: "7e74e2d88376bf641107fe22ea09ce68",
		},
		{
			name: "Two workflows",
			given: []Workflow{
				{
					Name:        "Workflow A",
					Description: "Test workflows A",
					Filename:    ".github/workflows/a.yml",
				},
				{
					Name:        "Workflow B",
					Description: "Test workflows B",
					Filename:    ".github/workflows/b.yml",
				},
			},
			expectedHash: "e6b8867c16847397e0eab0d694a39526",
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

			// Check # of items in the list
			if len(ws.Content.Items) != len(tt.given) {
				t.Errorf(errorf, "contents size mismatch", len(tt.given), len(ws.Content.Items))
			}

			// Check # of workflows
			if len(ws.Workflows) != len(tt.given) {
				t.Errorf(errorf, "workflows size mismatch", len(tt.given), len(ws.Workflows))
			}

			// Check Markdown
			if tt.expectedHash != helper.Hash(ws.Markdown()) {
				t.Errorf(errorf, "mismatch", tt.expectedHash, helper.Hash(ws.Markdown()))
				t.Error(ws.Markdown())
			}
		})
	}
}
