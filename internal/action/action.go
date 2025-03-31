package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nu12/action-docs/internal/markdown"
	"github.com/nu12/action-docs/internal/types"
	"github.com/nu12/go-logging"
	"gopkg.in/yaml.v3"
)

type Action struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Inputs      *types.InputMap  `yaml:"inputs"`
	Outputs     *types.OutputMap `yaml:"outputs"`
	Filename    string
}

func (a *Action) Markdown() string {
	inputs, outputs := a.getInputsOutputs()
	md := &markdown.Markdown{}
	md.Add(markdown.H1(a.Name)).
		Add(markdown.P(a.Description)).
		Add(markdown.H2("Usage example")).
		Add(markdown.Code(fmt.Sprintf("jobs:\n  job-name:\n    runs-on: <runner>\n    steps:\n    - uses: %s@main\n%s", filepath.Dir(a.Filename), inputs.ToString(8))))

	if len(*inputs) > 0 {
		md.Add(markdown.H2("Inputs"))

		tInputs := markdown.Table{
			Header: markdown.Header{"Name", "Description", "Required", "Default value"},
		}
		for name, input := range *inputs {
			tInputs.AddRow(markdown.Row{name, input.Description, strconv.FormatBool(input.Required), input.Default})
		}

		md.Add(tInputs.Sort(0))
	}

	if len(*outputs) > 0 {
		md.Add(markdown.H2("Outputs"))

		outputs := markdown.Table{
			Header: markdown.Header{"Name", "Description"},
		}
		for name, output := range *a.Outputs {
			outputs.AddRow(markdown.Row{name, output.Description})
		}

		md.Add(outputs.Sort(0))
	}

	return md.String()
}

func Parse(file string, log *logging.Log) *Action {
	a := &Action{
		Inputs:  &types.InputMap{},
		Outputs: &types.OutputMap{},
	}
	b, err := os.ReadFile(file)
	if err != nil {
		log.Warning(err.Error())
		return a
	}

	err = yaml.Unmarshal([]byte(b), a)
	if err != nil {
		log.Warning(err.Error())
	}
	a.Filename = file

	return a
}

func (a *Action) getInputs() *types.InputMap {
	if a.Inputs == nil {
		return &types.InputMap{}
	}
	a.Inputs.Sort()
	return a.Inputs
}

func (a *Action) getOutputs() *types.OutputMap {
	if a.Outputs == nil {
		return &types.OutputMap{}
	}
	a.Outputs.Sort()
	return a.Outputs
}

func (a *Action) getInputsOutputs() (*types.InputMap, *types.OutputMap) {
	return a.getInputs(), a.getOutputs()
}
