package workflow

import (
	"os"
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/go-logging"

	"gopkg.in/yaml.v3"
)

const errorf = "Error: %v. \nExpected: %v \nGot: %v"

func TestParse(t *testing.T) {
	tests := []struct {
		name                       string
		data                       string
		expectedFilename           string
		expectedIsReusableWorkflow bool
		expectedName               string
		expectedDescription        string
	}{
		{
			name: "Workflow call",
			data: `
name: 'Workflow name 1'
description: 'Workflow description 1'
on: 
  workflow_call:
    inputs: 
      in1: 
        description: 'Input1'
        required: true
      in2: 
        description: 'Input2'
        required: false
    outputs:
      out1:
        description: 'Output1'
        value: 'Hello'
    secrets:
      sec1:
        required: true
`,
			expectedFilename:           "call.yml",
			expectedIsReusableWorkflow: true,
			expectedName:               "Workflow name 1",
			expectedDescription:        "Workflow description 1",
		},
		{
			name: "Workflow dispatch",
			data: `
name: 'Workflow name 2'
description: 'Workflow description 2'
on: 
  workflow_dispatch:
    inputs: 
      in1: 
        description: 'Input1'
        type: choice
        default: 'one'
        options:
        - one
        - two
`,
			expectedFilename:           "dispatch.yml",
			expectedIsReusableWorkflow: false,
			expectedName:               "Workflow name 2",
			expectedDescription:        "Workflow description 2",
		},
		{
			name: "Workflow call with nil",
			data: `
name: 'Workflow name 3'
description: 'Workflow description 3'
on: 
  workflow_call: {}
`,
			expectedFilename:           "call.yml",
			expectedIsReusableWorkflow: true,
			expectedName:               "Workflow name 3",
			expectedDescription:        "Workflow description 3",
		},
		{
			name: "Workflow dispatch with nil",
			data: `
name: 'Workflow name 4'
description: 'Workflow description 4'
on: 
  workflow_dispatch: {}
`,
			expectedFilename:           "dispatch.yml",
			expectedIsReusableWorkflow: false,
			expectedName:               "Workflow name 4",
			expectedDescription:        "Workflow description 4",
		},
	}

	log := logging.NewLogger()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			absoluteFilename := dir + "/" + tt.expectedFilename
			err := os.WriteFile(absoluteFilename, []byte(tt.data), 0644)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			defer os.Remove(absoluteFilename)

			w := Parse(absoluteFilename, log)

			if w.IsReusableWorkflow() != tt.expectedIsReusableWorkflow {
				t.Errorf(errorf, "IsReusableWorkflow doesn't match", tt.expectedIsReusableWorkflow, w.IsReusableWorkflow())
			}

			if w.Name != tt.expectedName {
				t.Errorf(errorf, "Name doesn't match", tt.expectedName, w.Name)
			}

			if w.Description != tt.expectedDescription {
				t.Errorf(errorf, "Description doesn't match", tt.expectedDescription, w.Description)
			}

			if w.IsReusableWorkflow() {
				if w.On.WorkflowCall == nil {
					t.Errorf("WorkflowCall is nil")
				}
				if w.On.WorkflowCall.Inputs == nil {
					t.Errorf("Inputs is nil")
				}
				if w.On.WorkflowCall.Outputs == nil {
					t.Errorf("Outputs is nil")
				}
				if w.On.WorkflowCall.Secrets == nil {
					t.Errorf("Secrets is nil")
				}

			} else {
				if w.On.WorkflowDispatch == nil {
					t.Errorf("WorkflowDispatch is nil")
				}

				if w.On.WorkflowDispatch.Inputs == nil {
					t.Errorf("Inputs is nil")
				}
			}
		})
	}

}

func TestMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		filename string
		expected string
	}{
		{
			name: "Workflow call",
			data: `
name: 'Workflow name'
description: 'Workflow description'
on: 
  workflow_call:
    inputs: 
      in1: 
        description: 'Input1'
        required: true
      in2: 
        description: 'Input2'
        required: false
    outputs:
      out1:
        description: 'Output1'
        value: 'Hello'
    secrets:
      sec1:
        required: true
`,
			filename: ".github/workflows/call.yml",
			expected: "360289fb1c3e8e14b64cf0d592ebf21a",
		},
		{
			name: "Workflow dispatch",
			data: `
name: 'Workflow name'
description: 'Workflow description'
on: 
  workflow_dispatch:
    inputs: 
      in1: 
        description: 'Input1'
        type: choice
        default: 'one'
        options:
        - one
        - two
`,
			filename: ".github/workflows/dispatch.yml",
			expected: "ec84abaf87f6ec0ab0f0c3208e493229",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := Workflow{}
			err := yaml.Unmarshal([]byte(tt.data), &w)
			if err != nil {
				t.Errorf("error: %v", err)
			}
			w.Filename = tt.filename

			if tt.expected != helper.Hash(w.Markdown()) {
				t.Errorf(errorf, "Markdown doesn't match", tt.expected, helper.Hash(w.Markdown()))
				t.Error(w.Markdown())
			}
		})
	}
}
