package workflow

import (
	"os"
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/go-logging"

	"gopkg.in/yaml.v3"
)

var call = `
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
`

var dispatch = `
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
`

func TestReusableWorkflow(t *testing.T) {
	log := logging.NewLogger()
	dir := t.TempDir()
	valid := "valid.yml"

	err := os.WriteFile(dir+"/"+valid, []byte(call), 0644)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer os.Remove(dir + "/" + valid)

	w := Parse(dir+"/"+valid, log)

	if !w.IsReusableWorkflow() {
		t.Errorf("error: %s", "Should have workflow_call trigger")
	}

	if w.Name != "Workflow name" {
		t.Errorf("error: %s", "Name doesn't match")
	}

	inputs := w.getInputs()
	if len(*inputs) != 2 {
		t.Errorf("error: %s", "Should have 2 inputs")
	}

	if len(*w.On.WorkflowCall.Outputs) != 1 {
		t.Errorf("error: %s", "Should have 1 output")
	}

	if (*inputs)["in1"].Description != "Input1" {
		t.Errorf("error: %s", "Description for input1 doesn't match")
	}

	if !(*inputs)["in1"].Required {
		t.Errorf("error: %s", "Required for input1 doesn't match")
	}

	if (*inputs)["in2"].Description != "Input2" {
		t.Errorf("error: %s", "Description for input2 doesn't match")
	}

	if (*inputs)["in2"].Required {
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

func TestMarkdownCall(t *testing.T) {
	w := Workflow{}

	err := yaml.Unmarshal([]byte(call), &w)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	w.Filename = ".github/workflows/workflow.yml"

	result := w.Markdown()
	expected := "07c0de5551eea7025970cc8f3e78b564"

	if expected != helper.Hash(result) {
		t.Errorf("error: %s. Output is:\n%s\nCurrent Hash is: %s, expected hash is: %s", "Markdown doesn't match", result, helper.Hash(result), expected)
	}
}

func TestMarkdownDispatch(t *testing.T) {

	w := Workflow{}

	err := yaml.Unmarshal([]byte(dispatch), &w)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	w.Filename = ".github/workflows/workflow.yml"

	result := w.Markdown()
	expected := "fb9150a4dca978c8672289ad53192481"

	if expected != helper.Hash(result) {
		t.Errorf("error: %s. Output is:\n%s\nCurrent Hash is: %s, expected hash is: %s", "Markdown doesn't match", result, helper.Hash(result), expected)
	}
}

func TestListInputs(t *testing.T) {
	inputs := &map[string]Input{
		"in2": {Description: "Input2", Required: false},
		"in1": {Description: "Input1", Required: true},
	}
	result := listInputs(inputs, 2)
	expected := "  in1: \n  in2: \n"

	if result != expected {
		t.Errorf("error: %s. Output is:\n%s\nExpected output is:\n%s", "listInputs doesn't match", result, expected)
	}
}
