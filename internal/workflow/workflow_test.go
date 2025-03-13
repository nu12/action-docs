package workflow

import (
	"testing"

	"github.com/nu12/action-docs/internal/helper"

	"gopkg.in/yaml.v3"
)

var data = `
name: 'Workflow name'
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
`

func TestReusableWorkflow(t *testing.T) {

	w := Workflow{}

	err := yaml.Unmarshal([]byte(data), &w)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if !w.IsReusableWorkflow() {
		t.Errorf("error: %s", "Should have workflow_call trigger")
	}

	if w.Name != "Workflow name" {
		t.Errorf("error: %s", "Name doesn't match")
	}

	if len(*w.On.WorkflowCall.Inputs) != 2 {
		t.Errorf("error: %s", "Should have 2 inputs")
	}

	if len(*w.On.WorkflowCall.Outputs) != 1 {
		t.Errorf("error: %s", "Should have 1 output")
	}

	if (*w.On.WorkflowCall.Inputs)["in1"].Description != "Input1" {
		t.Errorf("error: %s", "Description for input1 doesn't match")
	}

	if !(*w.On.WorkflowCall.Inputs)["in1"].Required {
		t.Errorf("error: %s", "Required for input1 doesn't match")
	}

	if (*w.On.WorkflowCall.Inputs)["in2"].Description != "Input2" {
		t.Errorf("error: %s", "Description for input2 doesn't match")
	}

	if (*w.On.WorkflowCall.Inputs)["in2"].Required {
		t.Errorf("error: %s", "Required for input2 doesn't match")
	}

	if !(*w.On.WorkflowCall.Secrets)["sec1"].Required {
		t.Errorf("error: %s", "Required for secret1 doesn't match")
	}

	// d, err := yaml.Marshal(&w)
	// if err != nil {
	// 	t.Errorf("error 2: %v", err)
	// }
	// fmt.Printf("--- t dump:\n%s\n\n", string(d))
}

func TestMarkdown(t *testing.T) {
	w := Workflow{}

	err := yaml.Unmarshal([]byte(data), &w)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	result := w.Markdown()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	expected := "3c702ad36af5342a27057ce620ccd8d6"

	if expected != helper.Hash(result) {
		t.Errorf("error: %s. Output is:\n%s\nCurrent Hash is: %s, expected hash is: %s", "Markdown doesn't match", result, helper.Hash(result), expected)
	}
}
