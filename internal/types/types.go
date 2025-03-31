package types

import (
	"fmt"
	"sort"
	"strings"
)

type Input struct {
	Default     string `yaml:"default,omitempty"`
	Description string `yaml:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty"`
	Type        string `yaml:"type,omitempty"`
}

type Output struct {
	Description string `yaml:"description,omitempty"`
}

type Secret struct {
	Required bool `yaml:"required,omitempty"`
}

type InputMap map[string]Input
type OutputMap map[string]Output
type SecretMap map[string]Secret

func (left *InputMap) Equals(right *InputMap) bool {
	return left.ToString(2) == right.ToString(2)
}

func (left *OutputMap) Equals(right *OutputMap) bool {
	left.Sort()
	right.Sort()
	if len(*left) != len(*right) {
		return false
	}
	for name, leftItem := range *left {
		rightItem, ok := (*right)[name]
		if !ok {
			return false
		}
		if leftItem.Description != rightItem.Description {
			return false
		}
	}
	return true
}

func (left *SecretMap) Equals(right *SecretMap) bool {
	left.Sort()
	right.Sort()
	if len(*left) != len(*right) {
		return false
	}
	for name, leftItem := range *left {
		rightItem, ok := (*right)[name]
		if !ok {
			return false
		}
		if leftItem.Required != rightItem.Required {
			return false
		}
	}
	return true
}

func (im *InputMap) Sort() {
	sorted := Sort[Input](*im)
	*im = sorted
}

func (om *OutputMap) Sort() {
	sorted := Sort[Output](*om)
	*om = sorted
}

func (sm *SecretMap) Sort() {
	sorted := Sort[Secret](*sm)
	*sm = sorted
}

func Sort[M Input | Output | Secret](m map[string]M) map[string]M {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := map[string]M{}
	for _, k := range keys {
		sorted[k] = m[k]
	}
	return sorted
}

func (im *InputMap) ToString(spacing int) string {
	if im == nil {
		return ""
	}
	if len(*im) == 0 {
		return ""
	}

	var result = []string{}
	for name, item := range *im {
		result = append(result, fmt.Sprintf("%s%s: %s\n", strings.Repeat(" ", spacing), name, item.Default))
	}
	sort.Strings(result)
	return fmt.Sprintf("%swith:\n%s", strings.Repeat(" ", spacing-2), strings.Join(result, ""))
}
