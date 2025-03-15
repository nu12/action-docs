package workflow

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/nu12/action-docs/internal/markdown"
	"github.com/nu12/go-logging"
	"gopkg.in/yaml.v3"
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
	inputs := w.getInputs()
	md := &markdown.Markdown{}
	md.Add(markdown.H2(w.Name)).
		Add(markdown.P("File: " + w.Filename)).
		Add(markdown.P(w.Description))

	if w.IsReusableWorkflow() {
		md.Add(markdown.H3("Usage example")).
			Add(markdown.Code(fmt.Sprintf(`name: My workflow
on:
  push:
    branches:
    - main

jobs:
  my-job:
    uses: %s@main
    with: 
%s`, w.Filename, listInputs(inputs, 6))))
	}

	if w.On.WorkflowCall == nil {
		return md.String()
	}

	if inputs != nil {
		md.Add(markdown.H3("Inputs"))

		in := markdown.Table{
			Header: markdown.Header{"Name", "Type", "Description", "Required"},
		}
		for name, input := range *inputs {
			in.AddRow(markdown.Row{name, input.Type, input.Description, strconv.FormatBool(input.Required)})
		}

		md.Add(in.Sort(0))
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

func Parse(file string, log *logging.Log) *Workflow {
	w := &Workflow{}
	b, err := os.ReadFile(file)
	if err != nil {
		log.Warning(err.Error())
		return w
	}

	err = yaml.Unmarshal([]byte(b), w)
	if err != nil {
		log.Warning(err.Error())
	}
	w.Filename = file
	return w
}

func (w *Workflow) getInputs() *map[string]Input {
	if w.On.WorkflowCall.Inputs == nil {
		return &map[string]Input{}
	}

	keys := make([]string, 0, len(*w.On.WorkflowCall.Inputs))
	for k := range *w.On.WorkflowCall.Inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]Input)
	for _, k := range keys {
		sorted[k] = (*w.On.WorkflowCall.Inputs)[k]
	}
	return &sorted
}

func listInputs(inputs *map[string]Input, spacing int) string {
	if inputs == nil {
		return ""
	}

	var result = []string{}
	for name := range *inputs {
		result = append(result, fmt.Sprintf("%s%s: \n", strings.Repeat(" ", spacing), name))
	}
	sort.Strings(result)
	return strings.Join(result, "")
}
