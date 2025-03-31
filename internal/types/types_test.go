package types

import "testing"

const errorf = "Error: %v. \nExpected: %v \nGot: %v"

func TestInputMapToString(t *testing.T) {
	tests := []struct {
		name     string
		given    *InputMap
		expected string
	}{
		{
			name: "Two inputs",
			given: &InputMap{
				"in2": {Description: "Input2", Required: false},
				"in1": {Description: "Input1", Required: true},
			},
			expected: "with:\n  in1: \n  in2: \n",
		},
		{
			name: "Input without default value",
			given: &InputMap{
				"in2": {Description: "Input2", Required: false},
				"in1": {Description: "Input1", Required: true},
			},
			expected: "with:\n  in1: \n  in2: \n",
		},
		{
			name:     "No inputs",
			given:    &InputMap{},
			expected: "",
		},
		{
			name:     "Nil inputs",
			given:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.given.ToString(2)
			if got != tt.expected {
				t.Errorf(errorf, "InputMap to string doesn't match", tt.expected, got)
			}

		})
	}
}

func TestEquals4InputMap(t *testing.T) {
	left := &InputMap{
		"input1": {Description: "desc1", Required: true, Default: "default1"},
		"input2": {Description: "desc2", Required: false, Default: "default2"},
	}
	tests := []struct {
		name     string
		right    *InputMap
		expected bool
	}{
		{
			name: "Identical maps",
			right: &InputMap{
				"input1": {Description: "desc1", Required: true, Default: "default1"},
				"input2": {Description: "desc2", Required: false, Default: "default2"},
			},
			expected: true,
		},
		{
			name: "Unordered maps",
			right: &InputMap{
				"input2": {Description: "desc2", Required: false, Default: "default2"},
				"input1": {Description: "desc1", Required: true, Default: "default1"},
			},
			expected: true,
		},
		{
			name: "Different keys",
			right: &InputMap{
				"input1": {Description: "desc1", Required: true, Default: "default1"},
				"input3": {Description: "desc2", Required: false, Default: "default2"},
			},

			expected: false,
		},
		{
			name: "Different default values",
			right: &InputMap{
				"input1": {Description: "desc1", Required: true, Default: "default1"},
				"input2": {Description: "desc2", Required: false, Default: "default3"},
			},

			expected: false,
		},
		{
			name: "Different size",
			right: &InputMap{
				"input1": {Description: "desc1", Required: true, Default: "default1"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := left.Equals(tt.right); got != tt.expected {
				t.Errorf(errorf, "input is not equal", got, tt.expected)
			}
		})
	}

}

func TestEquals4OutputMap(t *testing.T) {
	left := &OutputMap{
		"input1": {Description: "desc1"},
		"input2": {Description: "desc2"},
	}
	tests := []struct {
		name     string
		right    *OutputMap
		expected bool
	}{
		{
			name: "Identical maps",
			right: &OutputMap{
				"input1": {Description: "desc1"},
				"input2": {Description: "desc2"},
			},
			expected: true,
		},
		{
			name: "Unordered maps",
			right: &OutputMap{
				"input2": {Description: "desc2"},
				"input1": {Description: "desc1"},
			},
			expected: true,
		},
		{
			name: "Different keys",
			right: &OutputMap{
				"input3": {Description: "desc2"},
				"input1": {Description: "desc1"},
			},
			expected: false,
		},
		{
			name: "Different description",
			right: &OutputMap{
				"input2": {Description: "desc3"},
				"input1": {Description: "desc1"},
			},
			expected: false,
		},
		{
			name: "Different size",
			right: &OutputMap{
				"input1": {Description: "desc1"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := left.Equals(tt.right); got != tt.expected {
				t.Errorf(errorf, "output is not equal", got, tt.expected)
			}
		})
	}

}

func TestEquals4SecretMap(t *testing.T) {
	left := &SecretMap{
		"secret1": {Required: true},
		"secret2": {Required: false},
	}
	tests := []struct {
		name     string
		right    *SecretMap
		expected bool
	}{
		{
			name: "Identical secrets",
			right: &SecretMap{
				"secret1": {Required: true},
				"secret2": {Required: false},
			},
			expected: true,
		},
		{
			name: "Unordered secrets",
			right: &SecretMap{
				"secret2": {Required: false},
				"secret1": {Required: true},
			},
			expected: true,
		},
		{
			name: "Different keys",
			right: &SecretMap{
				"secret3": {Required: false},
				"secret1": {Required: true},
			},
			expected: false,
		},
		{
			name: "Different size",
			right: &SecretMap{
				"secret1": {Required: true},
			},
			expected: false,
		},
		{
			name: "Different required",
			right: &SecretMap{
				"secret1": {Required: true},
				"secret2": {Required: true},
			},
			expected: false,
		},
		{
			name:     "Empty secret",
			right:    &SecretMap{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := left.Equals(tt.right); got != tt.expected {
				t.Errorf(errorf, "secret is not equal", got, tt.expected)
			}
		})
	}

}
