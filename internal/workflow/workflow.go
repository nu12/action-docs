package workflow

import (
	"strconv"

	"github.com/nu12/action-docs/internal/markdown"
)

type Workflow struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	On          struct {
		WorkflowCall *struct {
			Inputs  *map[string]Input  `yaml:"inputs"`
			Outputs *map[string]Output `yaml:"outputs"`
			Secrets *map[string]Secret `yaml:"secrets"`
		} `yaml:"workflow_call"`
	}

	Filename string
}

type Input struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Type        string `yaml:"type"`
}
type Output struct {
	Description string `yaml:"description"`
}
type Secret struct {
	Required bool `yaml:"required"`
}

func (w *Workflow) IsReusableWorkflow() bool {
	return w.On.WorkflowCall != nil
}

func (w *Workflow) Markdown() string {
	//md := markdown.H1(w.Name).String()
	md := &markdown.Markdown{}
	md.Add(markdown.H2(w.Name)).
		Add(markdown.P("File: " + w.Filename)).
		Add(markdown.P(w.Description))

	if w.IsReusableWorkflow() {
		md.Add(markdown.H3("Usage example")).
			Add(markdown.Code(`name: My workflow
on:
  push:
    branches:
    - main

jobs:
  my-job:
    uses: .github/workflows/<filename>.yml@main
    with:
      <list-of-inputs>`))
	}

	if w.On.WorkflowCall == nil {
		return md.String()
	}

	if w.On.WorkflowCall.Inputs != nil {
		md.Add(markdown.H3("Inputs"))

		inputs := markdown.Table{
			Header: markdown.Header{"Name", "Type", "Description", "Required"},
		}
		for name, input := range *w.On.WorkflowCall.Inputs {
			inputs.AddRow(markdown.Row{name, input.Type, input.Description, strconv.FormatBool(input.Required)})
		}

		md.Add(inputs.Sort(0))
	}

	if w.On.WorkflowCall.Outputs != nil {
		md.Add(markdown.H3("Outputs"))

		outputs := markdown.Table{
			Header: markdown.Header{"Name", "Description"},
		}
		for name, output := range *w.On.WorkflowCall.Outputs {
			outputs.AddRow(markdown.Row{name, output.Description})
		}

		md.Add(outputs.Sort(0))
	}

	if w.On.WorkflowCall.Secrets != nil {
		md.Add(markdown.H3("Secrets"))

		secrets := markdown.Table{
			Header: markdown.Header{"Name", "Required"},
		}
		for name, secret := range *w.On.WorkflowCall.Secrets {
			secrets.AddRow(markdown.Row{name, strconv.FormatBool(secret.Required)})
		}

		md.Add(secrets.Sort(0))
	}

	return md.String()
}
