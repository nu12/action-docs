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
	md := markdown.H1(a.Name).String()

	md += markdown.P(a.Description).String()

	md += markdown.H2("Usage example").String()

	md += markdown.Code(`jobs:
  job-name:
    runs-on: <runner>
    steps:
    - uses: path/to/action/folder@main
      with:
        <list of inputs>`).String()

	if a.Inputs != nil {
		md += markdown.H2("Inputs").String()

		inputs := markdown.Table{
			Header: markdown.Header{"Name", "Description", "Required", "Default value"},
		}
		for name, input := range *a.Inputs {
			inputs.AddRow(markdown.Row{name, input.Description, strconv.FormatBool(input.Required), input.Default})
		}

		md += inputs.Sort(0).String()
	}

	if a.Outputs != nil {
		md += markdown.H2("Outputs").String()

		outputs := markdown.Table{
			Header: markdown.Header{"Name", "Description"},
		}
		for name, output := range *a.Outputs {
			outputs.AddRow(markdown.Row{name, output.Description})
		}

		md += outputs.Sort(0).String()
	}

	return md
}
