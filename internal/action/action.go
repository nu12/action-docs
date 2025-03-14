package action

import (
	"strconv"

	"github.com/nu12/action-docs/internal/markdown"
)

type Action struct {
	Name        string             `yaml:"name"`
	Description string             `yaml:"description"`
	Inputs      *map[string]Input  `yaml:"inputs"`
	Outputs     *map[string]Output `yaml:"outputs"`
}

type Input struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
}

type Output struct {
	Description string `yaml:"description"`
}

func (a *Action) Markdown() string {
	md := &markdown.Markdown{}
	md.Add(markdown.H1(a.Name)).
		Add(markdown.P(a.Description)).
		Add(markdown.H2("Usage example")).
		Add(markdown.Code(`jobs:
  job-name:
    runs-on: <runner>
    steps:
    - uses: path/to/action/folder@main
      with:
        <list of inputs>`))

	if a.Inputs != nil {
		md.Add(markdown.H2("Inputs"))

		inputs := markdown.Table{
			Header: markdown.Header{"Name", "Description", "Required", "Default value"},
		}
		for name, input := range *a.Inputs {
			inputs.AddRow(markdown.Row{name, input.Description, strconv.FormatBool(input.Required), input.Default})
		}

		md.Add(inputs.Sort(0))
	}

	if a.Outputs != nil {
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
