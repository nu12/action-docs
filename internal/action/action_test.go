package action

import (
	"os"
	"testing"

	"github.com/nu12/action-docs/internal/helper"
	"github.com/nu12/go-logging"
	"gopkg.in/yaml.v3"
)

var data = `
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
`

func TestAction(t *testing.T) {
	log := logging.NewLogger()
	dir := t.TempDir()
	valid := "action.yml"

	err := os.WriteFile(dir+"/"+valid, []byte(data), 0644)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer os.Remove(dir + "/" + valid)

	a := Parse(dir+"/"+valid, log)

	if a.Name != "Composite action" {
		t.Errorf("error: %s", "Name doesn't match")
	}

	if a.Filename != dir+"/"+valid {
		t.Errorf("error: %s", "Filename doesn't match")
	}

	if a.Description != "Description of the action" {
		t.Errorf("error: %s", "Description doesn't match")
	}

	inputs := a.getInputs()
	if (*inputs)["in1"].Description != "Input1" {
		t.Errorf("error: %s", "Input 1 description doesn't match")
	}

	if !(*inputs)["in1"].Required {
		t.Errorf("error: %s", "Input 1 required doesn't match")
	}

	if (*inputs)["in2"].Description != "Input2" {
		t.Errorf("error: %s", "Input 2 description doesn't match")
	}

	if (*inputs)["in2"].Required {
		t.Errorf("error: %s", "Input 2 required doesn't match")
	}

	if (*inputs)["in3"].Description != "Input3" {
		t.Errorf("error: %s", "Input 3 description doesn't match")
	}

	if (*inputs)["in3"].Default != "default value" {
		t.Errorf("error: %s", "Input 3 required doesn't match")
	}

	if (*a.Outputs)["out1"].Description != "Output1" {
		t.Errorf("error: %s", "Output 1 description doesn't match")
	}
}

func TestMarkdown(t *testing.T) {
	a := Action{}
	a.Filename = "actions/test/action.yml"

	err := yaml.Unmarshal([]byte(data), &a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	result := a.Markdown()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	expected := "adc43796f07ed2b769847ffbeeb57280"

	if expected != helper.Hash(result) {
		t.Errorf("error: %s. Output is:\n%s\nCurrent Hash is: %s, expected hash is: %s", "Markdown doesn't match", result, helper.Hash(result), expected)
	}
}

func TestListInputs(t *testing.T) {
	inputs := map[string]Input{
		"in2": {Description: "Input2", Required: false, Default: "default2"},
		"in1": {Description: "Input1", Required: true, Default: "default1"},
	}

	result := listInputs(&inputs, 2)
	expected := "  in1: default1\n  in2: default2\n"

	if result != expected {
		t.Errorf("error: %s. Output is:\n%s\nExpected output is:\n%s", "listInputs doesn't match", result, expected)
	}
}
