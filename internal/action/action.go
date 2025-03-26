package action

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/nu12/action-docs/internal/markdown"
	"github.com/nu12/go-logging"
	"gopkg.in/yaml.v3"
)

type Action struct {
	Name        string             `yaml:"name"`
	Description string             `yaml:"description"`
	Inputs      *map[string]Input  `yaml:"inputs"`
	Outputs     *map[string]Output `yaml:"outputs"`
	Filename    string
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
	inputs := a.getInputs()
	md := &markdown.Markdown{}
	md.Add(markdown.H1(a.Name)).
		Add(markdown.P(a.Description)).
		Add(markdown.H2("Usage example")).
		Add(markdown.Code(fmt.Sprintf(`jobs:
  job-name:
    runs-on: <runner>
    steps:
    - uses: %s@main
%s`, filepath.Dir(a.Filename), listInputs(inputs, 8))))

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

	if a.Outputs != nil && len(*a.Outputs) > 0 {
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
	a := &Action{}
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

	if a.Inputs == nil {
		a.Inputs = &map[string]Input{}
	}
	if a.Outputs == nil {
		a.Outputs = &map[string]Output{}
	}

	return a
}

func (a *Action) getInputs() *map[string]Input {
	if a.Inputs == nil {
		return &map[string]Input{}
	}

	keys := make([]string, 0, len(*a.Inputs))
	for k := range *a.Inputs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make(map[string]Input)
	for _, k := range keys {
		sorted[k] = (*a.Inputs)[k]
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
	for name, item := range *inputs {
		result = append(result, fmt.Sprintf("%s%s: %s\n", strings.Repeat(" ", spacing), name, item.Default))
	}
	sort.Strings(result)
	return fmt.Sprintf("%swith:\n%s", strings.Repeat(" ", spacing-2), strings.Join(result, ""))
}
