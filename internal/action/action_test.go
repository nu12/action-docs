package action

import (
	"os"
	"sort"
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/types"
	"github.com/nu12/go-logging"
	"gopkg.in/yaml.v3"
)

const (
	errorf           = "Error: %v. \nExpected: %v \nGot: %v"
	expectedFilename = "action.yml"
)

var testData = []struct {
	name                string
	data                string
	filename            string
	expectedHash        string
	expectedName        string
	expectedDescription string
	expectedInputs      *types.InputMap
	expectedOutputs     *types.OutputMap
}{
	{
		name: "Complete action",
		data: `
name: 'Complete composite action'
description: 'Description of the complete action'
inputs:
  datain1:
    description: 'Input1 from data in'
    required: true
  datain2:
    description: 'Input2 from data in'
    required: false
  datain3:
    description: 'Input3 from data in'
    default: 'default value for datain3'
outputs:
  dataout1:
    description: 'Output from data out'
    value: 'Hello'
`,
		filename:            "actions/a/action.yml",
		expectedHash:        "bb5dd43b8ff38cb09999f920120be086",
		expectedName:        "Complete composite action",
		expectedDescription: "Description of the complete action",
		expectedInputs: &types.InputMap{
			"datain1": {Description: "Input1 from data in", Required: true},
			"datain2": {Description: "Input2 from data in", Required: false},
			"datain3": {Description: "Input3 from data in", Default: "default value for datain3"},
		},
		expectedOutputs: &types.OutputMap{
			"dataout1": {Description: "Output from data out"},
		},
	},
	{
		name: "Without outputs",
		data: `
name: 'Composite action without outputs'
description: 'Description of the action without outputs'
inputs:
  datain4:
    description: 'Input4 from data in'
    required: true
  datain5:
    description: 'Input5 from data in'
    required: false
  datain6:
    description: 'Input6 from data in'
    default: 'default value for datain6'
`,
		filename:            "actions/b/action.yml",
		expectedHash:        "733a912d26b8407601cfe669af4c0a05",
		expectedName:        "Composite action without outputs",
		expectedDescription: "Description of the action without outputs",
		expectedInputs: &types.InputMap{
			"datain4": {Description: "Input4 from data in", Required: true},
			"datain5": {Description: "Input5 from data in", Required: false},
			"datain6": {Description: "Input6 from data in", Default: "default value for datain6"},
		},
		expectedOutputs: &types.OutputMap{},
	},
	{
		name: "Without inputs",
		data: `
name: 'Composite action without inputs'
description: 'Description of the empty action'
`,
		filename:            "actions/c/action.yml",
		expectedHash:        "18136161a6d6a1957dd26daeda2b18d7",
		expectedName:        "Composite action without inputs",
		expectedDescription: "Description of the empty action",
		expectedInputs:      &types.InputMap{},
		expectedOutputs:     &types.OutputMap{},
	},
}

// Helper function to compare maps of inputs
func compareInputs(t *testing.T, expected, actual *types.InputMap) {
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
func compareOutputs(t *testing.T, expected, actual *types.OutputMap) {
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
func sortedKeys(m map[string]types.Input) []string {
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
			a.Filename = tt.filename
			err := yaml.Unmarshal([]byte(tt.data), &a)
			if err != nil {
				t.Errorf("error: %v", err)
			}

			if tt.expectedHash != helper.Hash(a.Markdown()) {
				t.Errorf(errorf, "Markdown doesn't match", tt.expectedHash, helper.Hash(a.Markdown()))
				t.Error(a.Markdown())
			}
		})
	}
}

func TestParse(t *testing.T) {
	log := logging.NewLogger()
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			absoluteFilename := dir + "/" + expectedFilename

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
