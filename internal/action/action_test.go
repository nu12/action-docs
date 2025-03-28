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

var testData = []struct {
	name                string
	data                string
	filename            string
	expectedHash        string
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
		filename:            "actions/test/action.yml",
		expectedHash:        "adc43796f07ed2b769847ffbeeb57280",
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
		filename:            "actions/test/action.yml",
		expectedHash:        "3c0052d79bab651a2844a74c343b340c",
		expectedName:        "Composite action",
		expectedDescription: "Description of the action",
		expectedFilename:    "action.yml",
		expectedInputs: &map[string]Input{
			"in1": {Description: "Input1", Required: true},
			"in2": {Description: "Input2", Required: false},
			"in3": {Description: "Input3", Default: "default value"},
		},
		expectedOutputs: &map[string]Output{},
	},
	{
		name: "Without inputs",
		data: `
name: 'Composite action'
description: 'Description of the empty action'
`,
		filename:            "actions/test/action.yml",
		expectedHash:        "afefd9dc26139ad4cbe32dcd3269a4aa",
		expectedName:        "Composite action",
		expectedDescription: "Description of the empty action",
		expectedFilename:    "action.yml",
		expectedInputs:      &map[string]Input{},
		expectedOutputs:     &map[string]Output{},
	},
}

// Helper function to compare maps of inputs
func compareInputs(t *testing.T, expected, actual *map[string]Input) {
	if actual == nil {
		t.Errorf("Inputs is nil")
	}
	if len(*expected) != len(*actual) {
		t.Errorf(errorf, "Inputs length doesn't match", len(*expected), len(*actual))
	}

	expectedKeys := sortedKeys(*expected)
	actualKeys := sortedKeys(*actual)

	for i := range expectedKeys {
		if expectedKeys[i] != actualKeys[i] {
			t.Errorf(errorf, "Input keys order doesn't match", expectedKeys[i], actualKeys[i])
		}

		exp := (*expected)[expectedKeys[i]]
		act := (*actual)[expectedKeys[i]]

		if exp.Description != act.Description {
			t.Errorf(errorf, "Input description doesn't match", exp.Description, act.Description)
		}
		if exp.Required != act.Required {
			t.Errorf(errorf, "Input required doesn't match", exp.Required, act.Required)
		}
		if exp.Default != act.Default {
			t.Errorf(errorf, "Input default doesn't match", exp.Default, act.Default)
		}
	}
}

// Helper function to compare maps of outputs
func compareOutputs(t *testing.T, expected, actual *map[string]Output) {
	if actual == nil {
		t.Errorf("Outputs is nil")
	}

	if len(*expected) != len(*actual) {
		t.Errorf(errorf, "Outputs length doesn't match", len(*expected), len(*actual))
	}

	for name, exp := range *expected {
		act := (*actual)[name]
		if exp.Description != act.Description {
			t.Errorf(errorf, "Output description doesn't match", exp.Description, act.Description)
		}
	}
}

// Helper function to get sorted keys of a map
func sortedKeys(m map[string]Input) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func TestMarkdown(t *testing.T) {
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			a := Action{}
			a.Filename = "actions/test/action.yml"
			err := yaml.Unmarshal([]byte(tt.data), &a)
			if err != nil {
				t.Errorf("error: %v", err)
			}

			if tt.expectedHash != helper.Hash(a.Markdown()) {
				t.Errorf(errorf, "Markdown doesn't match", tt.expectedHash, helper.Hash(a.Markdown()))
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

			compareInputs(t, tt.expectedInputs, inputs)
		})
	}
}

func TestParse(t *testing.T) {
	log := logging.NewLogger()
	for _, tt := range testData {
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

			compareInputs(t, tt.expectedInputs, a.Inputs)
			compareOutputs(t, tt.expectedOutputs, a.Outputs)
		})
	}
}
