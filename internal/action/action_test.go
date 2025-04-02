package action

import (
	"os"
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/types"
	"github.com/nu12/go-logging"
)

const (
	errorf   = "Error: %v. \nExpected: %v \nGot: %v"
	filename = "action.yml"
)

func TestAction(t *testing.T) {
	tests := []struct {
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

	log := logging.NewLogger()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmp := t.TempDir()
			tmpFile := tmp + "/" + filename
			// Write data to file
			err := os.WriteFile(tmpFile, []byte(tt.data), 0644)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			defer os.Remove(tmpFile)

			// Parse and adjust filename
			a := Parse(tmpFile, log)
			a.Filename = tt.filename

			// Check Name
			if a.Name != tt.expectedName {
				t.Errorf(errorf, "Name doesn't match", tt.expectedName, a.Name)
			}

			// Check Description
			if a.Description != tt.expectedDescription {
				t.Errorf(errorf, "Description doesn't match", tt.expectedDescription, a.Description)
			}

			// Check Filename
			if a.Filename != tt.filename {
				t.Errorf(errorf, "Filename doesn't match", tt.filename, a.Filename)
			}

			inputs, outputs := a.getInputsOutputs()

			// Check Inputs
			if !inputs.Equals(tt.expectedInputs) {
				t.Errorf(errorf, "Inputs don't match", tt.expectedInputs, inputs)
			}

			// Check Outputs
			if !outputs.Equals(tt.expectedOutputs) {
				t.Errorf(errorf, "Outputs don't match", tt.expectedOutputs, outputs)
			}

			// Markdown and check hash
			md := a.Markdown()
			if helper.Hash(md) != tt.expectedHash {
				t.Errorf(errorf, "Hash doesn't match", tt.expectedHash, helper.Hash(md))
			}
		})
	}
}
