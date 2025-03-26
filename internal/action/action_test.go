package action

import (
	"os"
	"sort"
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/go-logging"
	"gopkg.in/yaml.v3"
)

const errorf = "Error: %v. \nExpected: %v \nGot: %v"

func TestMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		filename string
		expected string
	}{
		{
			name: "Complete action",
			data: `
name: 'Composite action'
description: 'Description of the action'
inputs:
  in1:
    description: 'Input1'
    required: true
  in2:
    description: 'Input2'
    required: false
  in3:
    description: 'Input3'
    default: 'default value'
outputs:
  out1:
    description: 'Output1'
    value: 'Hello'
`,
			filename: "actions/test/action.yml",
			expected: "adc43796f07ed2b769847ffbeeb57280",
		},
		{
			name: "Without outputs",
			data: `
name: 'Composite action'
description: 'Description of the action'
inputs:
  in1:
    description: 'Input1'
    required: true
  in2:
    description: 'Input2'
    required: false
  in3:
    description: 'Input3'
    default: 'default value'
`,
			filename: "actions/test/action.yml",
			expected: "3c0052d79bab651a2844a74c343b340c",
		},
		{
			name: "Whithout inputs",
			data: `
name: 'Composite action'
description: 'Description of the action'
`,
			filename: "actions/test/action.yml",
			expected: "b3f3f6803829df63fb679c51b5a4ef70",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Action{}
			a.Filename = "actions/test/action.yml"
			err := yaml.Unmarshal([]byte(tt.data), &a)
			if err != nil {
				t.Errorf("error: %v", err)
			}

			if tt.expected != helper.Hash(a.Markdown()) {
				t.Errorf(errorf, "Markdown doesn't match", tt.expected, helper.Hash(a.Markdown()))
				t.Errorf(a.Markdown())
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
				"in2": {Description: "Input2", Required: false, Default: "default2"},
				"in1": {Description: "Input1", Required: true, Default: "default1"},
			},
			expected: "with:\n  in1: default1\n  in2: default2\n",
		},
		{
			name: "Input without default value",
			given: map[string]Input{
				"in2": {Description: "Input2", Required: false, Default: "default2"},
				"in1": {Description: "Input1", Required: true},
			},
			expected: "with:\n  in1: \n  in2: default2\n",
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
		given          Action
		expectedInputs *map[string]Input
	}{
		{
			name: "Action with inputs",
			given: Action{
				Inputs: &map[string]Input{
					"in3": {Description: "Input3", Default: "default value"},
					"in2": {Description: "Input2", Required: false},
					"in1": {Description: "Input1", Required: true},
				},
			},
			expectedInputs: &map[string]Input{
				"in1": {Description: "Input1", Required: true},
				"in2": {Description: "Input2", Required: false},
				"in3": {Description: "Input3", Default: "default value"},
			},
		},
		{
			name:           "Action with nil inputs",
			given:          Action{},
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

				if (*inputs)[expectedKeys[i]].Default != (*tt.expectedInputs)[expectedKeys[i]].Default {
					t.Errorf(errorf, "Input default doesn't match", (*tt.expectedInputs)[expectedKeys[i]].Default, (*inputs)[expectedKeys[i]].Default)
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name                string
		data                string
		expectedName        string
		expectedDescription string
		expectedFilename    string
		expectedInputs      *map[string]Input
		expectedOutputs     *map[string]Output
	}{
		{
			name: "Complete action",
			data: `
name: 'Composite action'
description: 'Description of the action'
inputs:
  in1:
    description: 'Input1'
    required: true
  in2:
    description: 'Input2'
    required: false
  in3:
    description: 'Input3'
    default: 'default value'
outputs:
  out1:
    description: 'Output1'
    value: 'Hello'
`,
			expectedName:        "Composite action",
			expectedDescription: "Description of the action",
			expectedFilename:    "action.yml",
			expectedInputs: &map[string]Input{
				"in1": {Description: "Input1", Required: true},
				"in2": {Description: "Input2", Required: false},
				"in3": {Description: "Input3", Default: "default value"},
			},
			expectedOutputs: &map[string]Output{
				"out1": {Description: "Output1"},
			},
		},
		{
			name: "Empty action",
			data: `
name: 'Composite action'
description: 'Description of the empty action'
`,
			expectedName:        "Composite action",
			expectedDescription: "Description of the empty action",
			expectedFilename:    "action.yml",
			expectedInputs:      &map[string]Input{},
			expectedOutputs:     &map[string]Output{},
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

			a := Parse(absoluteFilename, log)

			if a.Name != tt.expectedName {
				t.Errorf(errorf, "Name doesn't match", tt.expectedName, a.Name)
			}

			if a.Filename != absoluteFilename {
				t.Errorf(errorf, "Filename doesn't match", absoluteFilename, a.Filename)
			}

			if a.Description != tt.expectedDescription {
				t.Errorf(errorf, "Description doesn't match", tt.expectedDescription, a.Description)
			}

			if a.Inputs == nil {
				t.Errorf("Inputs is nil")
			}

			if a.Outputs == nil {
				t.Errorf("Outputs is nil")
			}

			if len(*a.Inputs) != len(*tt.expectedInputs) {
				t.Errorf(errorf, "Inputs length doesn't match", len(*tt.expectedInputs), len(*a.Inputs))
			}

			if len(*a.Outputs) != len(*tt.expectedOutputs) {
				t.Errorf(errorf, "Outputs length doesn't match", len(*tt.expectedOutputs), len(*a.Outputs))
			}

			for name, input := range *tt.expectedInputs {
				if input.Description != (*a.Inputs)[name].Description {
					t.Errorf(errorf, "Input description doesn't match", input.Description, (*a.Inputs)[name].Description)
				}

				if input.Required != (*a.Inputs)[name].Required {
					t.Errorf(errorf, "Input required doesn't match", input.Required, (*a.Inputs)[name].Required)
				}

				if input.Default != (*a.Inputs)[name].Default {
					t.Errorf(errorf, "Input default doesn't match", input.Default, (*a.Inputs)[name].Default)
				}
			}

			for name, output := range *tt.expectedOutputs {
				if output.Description != (*a.Outputs)[name].Description {
					t.Errorf(errorf, "Input description doesn't match", output.Description, (*a.Outputs)[name].Description)
				}
			}
		})
	}
}
