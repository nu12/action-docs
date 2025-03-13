package workflow

import (
	"strconv"

	"github.com/nu12/action-docs/internal/markdown"
)

type Workflow struct {
	Name string
	On   struct {
		WorkflowCall *struct {
			Inputs  *map[string]Input  `yaml:"inputs"`
			Outputs *map[string]Output `yaml:"outputs"`
			Secrets *map[string]Secret `yaml:"secrets"`
		} `yaml:"workflow_call"`
	}
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
	md := markdown.H1(w.Name).String()

	if w.IsReusableWorkflow() {
		md += markdown.P("Reusable workflow").String()

		md += markdown.H2("Usage example").String()

		md += markdown.Code(`name: My workflow
on:
  push:
    branches:
    - main

jobs:
  my-job:
    uses: .github/workflows/<filename>.yml@main
    with:
      <list-of-inputs>`).String()
	}

	if w.On.WorkflowCall.Inputs != nil {
		md += markdown.H2("Inputs").String()

		inputs := markdown.Table{
			Header: markdown.Header{"Name", "Type", "Description", "Required"},
		}
		for name, input := range *w.On.WorkflowCall.Inputs {
			inputs.AddRow(markdown.Row{name, input.Type, input.Description, strconv.FormatBool(input.Required)})
		}

		md += inputs.Sort(0).String()
	}

	if w.On.WorkflowCall.Outputs != nil {
		md += markdown.H2("Outputs").String()

		outputs := markdown.Table{
			Header: markdown.Header{"Name", "Description"},
		}
		for name, output := range *w.On.WorkflowCall.Outputs {
			outputs.AddRow(markdown.Row{name, output.Description})
		}

		md += outputs.Sort(0).String()
	}

	if w.On.WorkflowCall.Secrets != nil {
		md += markdown.H2("Secrets").String()

		secrets := markdown.Table{
			Header: markdown.Header{"Name", "Required"},
		}
		for name, secret := range *w.On.WorkflowCall.Secrets {
			secrets.AddRow(markdown.Row{name, strconv.FormatBool(secret.Required)})
		}

		md += secrets.Sort(0).String()
	}

	return md
}
