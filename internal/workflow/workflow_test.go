package workflow

import (
	"os"
	"sort"
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
				t.Errorf(w.Markdown())
			}
		})
	}
}

func TestListInputs(t *testing.T) {
	tests := []struct {
		name     string
		given    map[string]Input
		expected string
	}{
		{
			name: "Two inputs",
			given: map[string]Input{
				"in2": {Description: "Input2", Required: false},
				"in1": {Description: "Input1", Required: true},
			},
			expected: "with:\n  in1: \n  in2: \n",
		},
		{
			name: "Input without default value",
			given: map[string]Input{
				"in2": {Description: "Input2", Required: false},
				"in1": {Description: "Input1", Required: true},
			},
			expected: "with:\n  in1: \n  in2: \n",
		},
		{
			name:     "No inputs",
			given:    map[string]Input{},
			expected: "",
		},
		{
			name:     "Nil inputs",
			given:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listInputs(&tt.given, 2)
			if got != tt.expected {
				t.Errorf(errorf, "List inputs doesn't match", tt.expected, got)
			}

		})
	}
}

func TestGetInputs(t *testing.T) {
	tests := []struct {
		name           string
		given          Workflow
		expectedInputs *map[string]Input
	}{
		{
			name: "Workflow dispatch with inputs",
			given: Workflow{
				Name:        "A",
				Description: "Test workflows",
				Filename:    ".github/workflows/a.yml",
				On: struct {
					WorkflowCall *struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					} "yaml:\"workflow_call\""
					WorkflowDispatch *struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					} "yaml:\"workflow_dispatch\""
				}{
					WorkflowDispatch: &struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					}{
						Inputs: &map[string]Input{
							"in2": {Description: "Input2", Required: false},
							"in1": {Description: "Input1", Required: true},
						},
					},
				},
			},
			expectedInputs: &map[string]Input{
				"in1": {Description: "Input1", Required: true},
				"in2": {Description: "Input2", Required: false},
			},
		},
		{
			name: "Workflow dispatch without inputs",
			given: Workflow{
				Name:        "A",
				Description: "Test workflows",
				Filename:    ".github/workflows/a.yml",
				On: struct {
					WorkflowCall *struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					} "yaml:\"workflow_call\""
					WorkflowDispatch *struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					} "yaml:\"workflow_dispatch\""
				}{
					WorkflowDispatch: &struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					}{
						Inputs: &map[string]Input{},
					},
				},
			},
			expectedInputs: &map[string]Input{},
		},
		{
			name: "Workflow dispatch with nil inputs",
			given: Workflow{
				Name:        "A",
				Description: "Test workflows",
				Filename:    ".github/workflows/a.yml",
				On: struct {
					WorkflowCall *struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					} "yaml:\"workflow_call\""
					WorkflowDispatch *struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					} "yaml:\"workflow_dispatch\""
				}{
					WorkflowDispatch: &struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					}{
						Inputs: nil,
					},
				},
			},
			expectedInputs: &map[string]Input{},
		},
		{
			name: "Workflow call with inputs",
			given: Workflow{
				Name:        "A",
				Description: "Test workflows",
				Filename:    ".github/workflows/a.yml",
				On: struct {
					WorkflowCall *struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					} "yaml:\"workflow_call\""
					WorkflowDispatch *struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					} "yaml:\"workflow_dispatch\""
				}{
					WorkflowCall: &struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					}{
						Inputs: &map[string]Input{
							"in2": {Description: "Input2", Required: false},
							"in1": {Description: "Input1", Required: true},
						},
						Outputs: &map[string]Output{},
						Secrets: &map[string]Secret{},
					},
				},
			},
			expectedInputs: &map[string]Input{
				"in1": {Description: "Input1", Required: true},
				"in2": {Description: "Input2", Required: false},
			},
		},
		{
			name: "Workflow call without inputs",
			given: Workflow{
				Name:        "A",
				Description: "Test workflows",
				Filename:    ".github/workflows/a.yml",
				On: struct {
					WorkflowCall *struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					} "yaml:\"workflow_call\""
					WorkflowDispatch *struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					} "yaml:\"workflow_dispatch\""
				}{
					WorkflowCall: &struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					}{
						Inputs:  &map[string]Input{},
						Outputs: &map[string]Output{},
						Secrets: &map[string]Secret{},
					},
				},
			},
			expectedInputs: &map[string]Input{},
		},
		{
			name: "Workflow call with nil inputs",
			given: Workflow{
				Name:        "A",
				Description: "Test workflows",
				Filename:    ".github/workflows/a.yml",
				On: struct {
					WorkflowCall *struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					} "yaml:\"workflow_call\""
					WorkflowDispatch *struct {
						Inputs *map[string]Input "yaml:\"inputs\""
					} "yaml:\"workflow_dispatch\""
				}{
					WorkflowCall: &struct {
						Inputs  *map[string]Input  "yaml:\"inputs\""
						Outputs *map[string]Output "yaml:\"outputs\""
						Secrets *map[string]Secret "yaml:\"secrets\""
					}{
						Inputs:  nil,
						Outputs: &map[string]Output{},
						Secrets: &map[string]Secret{},
					},
				},
			},
			expectedInputs: &map[string]Input{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs := tt.given.getInputs()

			if inputs == nil {
				t.Errorf("Inputs is nil")
			}

			if len(*inputs) != len(*tt.expectedInputs) {
				t.Errorf(errorf, "Inputs length doesn't match", len(*tt.expectedInputs), len(*inputs))
			}

			var expectedKeys = make([]string, 0, len(*tt.expectedInputs))
			for key, _ := range *tt.expectedInputs {
				expectedKeys = append(expectedKeys, key)
			}
			sort.Strings(expectedKeys)

			var gotKeys = make([]string, 0, len(*inputs))
			for key, _ := range *inputs {
				gotKeys = append(gotKeys, key)
			}

			for i := range expectedKeys {
				if expectedKeys[i] != gotKeys[i] {
					t.Errorf(errorf, "Input keys order doesn't match", expectedKeys[i], gotKeys[i])
				}

				if (*inputs)[expectedKeys[i]].Description != (*tt.expectedInputs)[expectedKeys[i]].Description {
					t.Errorf(errorf, "Input description doesn't match", (*tt.expectedInputs)[expectedKeys[i]].Description, (*inputs)[expectedKeys[i]].Description)
				}

				if (*inputs)[expectedKeys[i]].Required != (*tt.expectedInputs)[expectedKeys[i]].Required {
					t.Errorf(errorf, "Input required doesn't match", (*tt.expectedInputs)[expectedKeys[i]].Required, (*inputs)[expectedKeys[i]].Required)
				}
			}
		})
	}
}
