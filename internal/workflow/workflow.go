package workflow

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nu12/action-docs/internal/markdown"
	"github.com/nu12/action-docs/internal/types"
	"github.com/nu12/go-logging"
	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	On          struct {
		WorkflowCall *struct {
			Inputs  *types.InputMap  `yaml:"inputs"`
			Outputs *types.OutputMap `yaml:"outputs"`
			Secrets *types.SecretMap `yaml:"secrets"`
		} `yaml:"workflow_call"`
		WorkflowDispatch *struct {
			Inputs *types.InputMap `yaml:"inputs"`
		} `yaml:"workflow_dispatch"`
	}
	Filename           string
	IsReusableWorkflow bool
}

func (w *Workflow) Markdown() string {
	inputs, outputs, secrets := w.getInputsOutputsSecrets()
	md := &markdown.Markdown{}
	md.Add(markdown.H2(w.Name)).
		Add(markdown.P("File: " + w.Filename)).
		Add(markdown.P(w.Description))

	if w.IsReusableWorkflow {
		md.Add(markdown.H3("Usage example")).
			Add(markdown.Code(fmt.Sprintf("name: My workflow\non:\n  push:\n    branches:\n    - main\n\njobs:\n  my-job:\n    uses: %s@main\n%s", w.Filename, inputs.ToString(6))))
	}

	if len(*inputs) > 0 {
		md.Add(markdown.H3("Inputs"))

		if w.IsReusableWorkflow {
			tInputs := markdown.Table{
				Header: markdown.Header{"Name", "Type", "Description", "Required"},
			}
			for name, input := range *inputs {
				tInputs.AddRow(markdown.Row{name, input.Type, input.Description, strconv.FormatBool(input.Required)})
			}

			md.Add(tInputs.Sort(0))
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

	if len(*outputs) > 0 {
		md.Add(markdown.H3("Outputs"))

		tOutputs := markdown.Table{
			Header: markdown.Header{"Name", "Description"},
		}
		for name, output := range *w.On.WorkflowCall.Outputs {
			tOutputs.AddRow(markdown.Row{name, output.Description})
		}

		md.Add(tOutputs.Sort(0))
	}

	if len(*secrets) > 0 {
		md.Add(markdown.H3("Secrets"))

		tSecrets := markdown.Table{
			Header: markdown.Header{"Name", "Required"},
		}
		for name, secret := range *w.On.WorkflowCall.Secrets {
			tSecrets.AddRow(markdown.Row{name, strconv.FormatBool(secret.Required)})
		}

		md.Add(tSecrets.Sort(0))
	}

	return md.String()
}

func Parse(file string, log *logging.Log) *Workflow {
	w := &Workflow{
		On: struct {
			WorkflowCall *struct {
				Inputs  *types.InputMap  `yaml:"inputs"`
				Outputs *types.OutputMap `yaml:"outputs"`
				Secrets *types.SecretMap `yaml:"secrets"`
			} `yaml:"workflow_call"`
			WorkflowDispatch *struct {
				Inputs *types.InputMap `yaml:"inputs"`
			} `yaml:"workflow_dispatch"`
		}{
			WorkflowCall: &struct {
				Inputs  *types.InputMap  `yaml:"inputs"`
				Outputs *types.OutputMap `yaml:"outputs"`
				Secrets *types.SecretMap `yaml:"secrets"`
			}{
				Inputs:  &types.InputMap{},
				Outputs: &types.OutputMap{},
				Secrets: &types.SecretMap{},
			},
			WorkflowDispatch: &struct {
				Inputs *types.InputMap `yaml:"inputs"`
			}{
				Inputs: &types.InputMap{},
			},
		},
		IsReusableWorkflow: false,
	}

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

	if strings.Contains(string(b), "workflow_call") {
		w.IsReusableWorkflow = true
	}
	return w
}

func (w *Workflow) getInputs() *types.InputMap {
	if w.IsReusableWorkflow {
		if w.On.WorkflowCall.Inputs == nil {
			return &types.InputMap{}
		}
		w.On.WorkflowCall.Inputs.Sort()
		return w.On.WorkflowCall.Inputs
	}
	if w.On.WorkflowDispatch == nil {
		return &types.InputMap{}
	}
	if w.On.WorkflowDispatch.Inputs == nil {
		return &types.InputMap{}
	}
	w.On.WorkflowDispatch.Inputs.Sort()
	return w.On.WorkflowDispatch.Inputs
}

func (w *Workflow) getOutputs() *types.OutputMap {
	if w.IsReusableWorkflow {
		if w.On.WorkflowCall.Outputs == nil {
			return &types.OutputMap{}
		}
		w.On.WorkflowCall.Outputs.Sort()
		return w.On.WorkflowCall.Outputs
	}
	return &types.OutputMap{}
}

func (w *Workflow) getSecrets() *types.SecretMap {
	if w.IsReusableWorkflow {
		if w.On.WorkflowCall.Secrets == nil {
			return &types.SecretMap{}
		}
		w.On.WorkflowCall.Secrets.Sort()
		return w.On.WorkflowCall.Secrets
	}
	return &types.SecretMap{}
}

func (w *Workflow) getInputsOutputsSecrets() (*types.InputMap, *types.OutputMap, *types.SecretMap) {
	return w.getInputs(), w.getOutputs(), w.getSecrets()
}
