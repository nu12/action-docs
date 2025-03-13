package action

import (
	"testing"

	"github.com/nu12/action-docs/internal/helper"
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
	a := Action{}

	err := yaml.Unmarshal([]byte(data), &a)
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if a.Name != "Composite action" {
		t.Errorf("error: %s", "Name doesn't match")
	}

	if a.Description != "Description of the action" {
		t.Errorf("error: %s", "Description doesn't match")
	}

	if (*a.Inputs)["in1"].Description != "Input1" {
		t.Errorf("error: %s", "Input 1 description doesn't match")
	}

	if !(*a.Inputs)["in1"].Required {
		t.Errorf("error: %s", "Input 1 required doesn't match")
	}

	if (*a.Inputs)["in2"].Description != "Input2" {
		t.Errorf("error: %s", "Input 2 description doesn't match")
	}

	if (*a.Inputs)["in2"].Required {
		t.Errorf("error: %s", "Input 2 required doesn't match")
	}

	if (*a.Inputs)["in3"].Description != "Input3" {
		t.Errorf("error: %s", "Input 3 description doesn't match")
	}

	if (*a.Inputs)["in3"].Default != "default value" {
		t.Errorf("error: %s", "Input 3 required doesn't match")
	}

	if (*a.Outputs)["out1"].Description != "Output1" {
		t.Errorf("error: %s", "Output 1 description doesn't match")
	}
}

func TestMarkdown(t *testing.T) {
	a := Action{}

	err := yaml.Unmarshal([]byte(data), &a)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	result := a.Markdown()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	expected := "3f756d4abfc69e51c09b2532aef88bb5"

	if expected != helper.Hash(result) {
		t.Errorf("error: %s. Output is:\n%s\nCurrent Hash is: %s, expected hash is: %s", "Markdown doesn't match", result, helper.Hash(result), expected)
	}
}
