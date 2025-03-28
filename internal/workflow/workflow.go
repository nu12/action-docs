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
		WorkflowDispatch *struct {
			Inputs *map[string]Input `yaml:"inputs"`
		} `yaml:"workflow_dispatch"`
	}

	Filename           string
	isReusableWorkflow bool
}

type Input struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Type        string `yaml:"type"`
	Default     string `yaml:"default"`
}
type Output struct {
	Description string `yaml:"description"`
}
type Secret struct {
	Required bool `yaml:"required"`
}

func (w *Workflow) IsReusableWorkflow() bool {
	// In case the object is initialized without the Parse function
	if w.On.WorkflowCall == nil {
		return false
	}
	if w.On.WorkflowDispatch == nil {
		return true
	}

	// In case the object is initialized with the Parse function
	return w.isReusableWorkflow
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
%s`, w.Filename, listInputs(inputs, 6))))
	}

	if len(*inputs) > 0 {
		md.Add(markdown.H3("Inputs"))

		if w.IsReusableWorkflow() {
			in := markdown.Table{
				Header: markdown.Header{"Name", "Type", "Description", "Required"},
			}
			for name, input := range *inputs {
				in.AddRow(markdown.Row{name, input.Type, input.Description, strconv.FormatBool(input.Required)})
			}

			md.Add(in.Sort(0))
		} else {
			in := markdown.Table{
				Header: markdown.Header{"Name", "Type", "Description", "Default"},
			}
			for name, input := range *inputs {
				in.AddRow(markdown.Row{name, input.Type, input.Description, input.Default})
			}

			md.Add(in.Sort(0))
		}
	}

	if w.On.WorkflowCall == nil {
		return md.String()
	}

	if w.On.WorkflowCall.Outputs != nil && len(*w.On.WorkflowCall.Outputs) > 0 {
		md.Add(markdown.H3("Outputs"))

		outputs := markdown.Table{
			Header: markdown.Header{"Name", "Description"},
		}
		for name, output := range *w.On.WorkflowCall.Outputs {
			outputs.AddRow(markdown.Row{name, output.Description})
		}

		md.Add(outputs.Sort(0))
	}

	if w.On.WorkflowCall.Secrets != nil && len(*w.On.WorkflowCall.Secrets) > 0 {
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

	if w.On.WorkflowCall == nil {
		w.isReusableWorkflow = false
		w.On.WorkflowCall = &struct {
			Inputs  *map[string]Input  `yaml:"inputs"`
			Outputs *map[string]Output `yaml:"outputs"`
			Secrets *map[string]Secret `yaml:"secrets"`
		}{
			Inputs:  &map[string]Input{},
			Outputs: &map[string]Output{},
			Secrets: &map[string]Secret{},
		}
	}

	if w.On.WorkflowDispatch == nil {
		w.isReusableWorkflow = true
		w.On.WorkflowDispatch = &struct {
			Inputs *map[string]Input `yaml:"inputs"`
		}{
			Inputs: &map[string]Input{},
		}
	}

	if w.On.WorkflowCall.Inputs == nil {
		w.On.WorkflowCall.Inputs = &map[string]Input{}
	}
	if w.On.WorkflowCall.Outputs == nil {
		w.On.WorkflowCall.Outputs = &map[string]Output{}
	}
	if w.On.WorkflowCall.Secrets == nil {
		w.On.WorkflowCall.Secrets = &map[string]Secret{}
	}
	if w.On.WorkflowDispatch.Inputs == nil {
		w.On.WorkflowDispatch.Inputs = &map[string]Input{}
	}

	return w
}

func (w *Workflow) getInputs() *map[string]Input {
	sorted := make(map[string]Input)
	var keys = []string{}
	if w.On.WorkflowCall != nil {
		if w.On.WorkflowCall.Inputs != nil {
			for k := range *w.On.WorkflowCall.Inputs {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				sorted[k] = (*w.On.WorkflowCall.Inputs)[k]
			}
		}
	} else if w.On.WorkflowDispatch != nil {
		if w.On.WorkflowDispatch.Inputs != nil {
			for k := range *w.On.WorkflowDispatch.Inputs {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				sorted[k] = (*w.On.WorkflowDispatch.Inputs)[k]
			}
		}
	}
	return &sorted
}

func listInputs(inputs *map[string]Input, spacing int) string {
	if inputs == nil {
		return ""
	}
	if len(*inputs) == 0 {
		return ""
	}

	var result = []string{}
	for name := range *inputs {
		result = append(result, fmt.Sprintf("%s%s: \n", strings.Repeat(" ", spacing), name))
	}
	sort.Strings(result)
	return fmt.Sprintf("%swith:\n%s", strings.Repeat(" ", spacing-2), strings.Join(result, ""))
}
