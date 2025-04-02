package workflow

import (
	"os"
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/action-docs/internal/types"
	"github.com/nu12/go-logging"
)

const errorf = "Error: %v. \nExpected: %v \nGot: %v"

func TestWorkflow(t *testing.T) {
	// Test data
	tests := []struct {
		name                       string
		data                       string
		expectedFilename           string
		expectedIsReusableWorkflow bool
		expectedName               string
		expectedDescription        string
		expectedHash               string
		expectedInputs             *types.InputMap
		expectedOutputs            *types.OutputMap
		expectedSecrets            *types.SecretMap
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
			expectedHash:               "0e6209aadcd84a36add6dac6aaacc09f",
			expectedInputs: &types.InputMap{
				"in1": {Description: "Input1", Required: true},
				"in2": {Description: "Input2", Required: false},
			},
			expectedOutputs: &types.OutputMap{
				"out1": {Description: "Output1"},
			},
			expectedSecrets: &types.SecretMap{
				"sec1": {Required: true},
			},
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
			expectedHash:               "8721bc3d2f959c9a5fe825d60eb70514",
			expectedInputs: &types.InputMap{
				"in1": {Description: "Input1", Type: "choice", Default: "one"},
			},
			expectedOutputs: &types.OutputMap{},
			expectedSecrets: &types.SecretMap{},
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
			expectedHash:               "696b5d1e79c3dbff8476ec1dc47643e6",
			expectedInputs:             &types.InputMap{},
			expectedOutputs:            &types.OutputMap{},
			expectedSecrets:            &types.SecretMap{},
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
			expectedHash:               "a90ad004f3a963f718a317ba6dd55f0f",
			expectedInputs:             &types.InputMap{},
			expectedOutputs:            &types.OutputMap{},
			expectedSecrets:            &types.SecretMap{},
		},
	}

	log := logging.NewLogger()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tmpFile := dir + "/" + tt.expectedFilename
			err := os.WriteFile(tmpFile, []byte(tt.data), 0644)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			defer os.Remove(tmpFile)

			// Parse
			w := Parse(tmpFile, log)
			w.Filename = tt.expectedFilename
			// Check name
			if w.Name != tt.expectedName {
				t.Errorf(errorf, "Name doesn't match", tt.expectedName, w.Name)
			}

			// Check filename
			if w.Filename != tt.expectedFilename {
				t.Errorf(errorf, "Filename doesn't match", tt.expectedFilename, w.Filename)
			}

			// Check description
			if w.Description != tt.expectedDescription {
				t.Errorf(errorf, "Description doesn't match", tt.expectedDescription, w.Description)
			}

			// Check is reusable workflow
			if w.IsReusableWorkflow != tt.expectedIsReusableWorkflow {
				t.Errorf(errorf, "IsReusableWorkflow doesn't match", tt.expectedIsReusableWorkflow, w.IsReusableWorkflow)
			}

			inputs, outputs, secrets := w.getInputsOutputsSecrets()

			// Check inputs
			if !inputs.Equals(tt.expectedInputs) {
				t.Errorf(errorf, "Inputs don't match", tt.expectedInputs, inputs)
			}

			// Check outputs
			if !outputs.Equals(tt.expectedOutputs) {
				t.Errorf(errorf, "Outputs don't match", tt.expectedOutputs, outputs)
			}

			// Check secrets
			if !secrets.Equals(tt.expectedSecrets) {
				t.Errorf(errorf, "Secrets don't match", tt.expectedSecrets, secrets)
			}

			// Check Markdown
			md := w.Markdown()
			if tt.expectedHash != helper.Hash(md) {
				t.Errorf(errorf, "Markdown doesn't match", tt.expectedHash, helper.Hash(md))
				t.Error(md)
			}

		})
	}

}
