package markdown

import "testing"

func TestHyperlink(t *testing.T) {
	tests := []struct {
		name      string
		hyperlink Hyperlink
		expected  string
	}{
		{
			name: "Valid hyperlink",
			hyperlink: Hyperlink{
				URL:  "https://example.com",
				Text: "Example",
			},
			expected: "[Example](https://example.com)",
		},
		{
			name: "Empty URL",
			hyperlink: Hyperlink{
				URL:  "",
				Text: "No URL",
			},
			expected: "[No URL]()",
		},
		{
			name: "Empty Text",
			hyperlink: Hyperlink{
				URL:  "https://example.com",
				Text: "",
			},
			expected: "[](https://example.com)",
		},
		{
			name: "Empty URL and Text",
			hyperlink: Hyperlink{
				URL:  "",
				Text: "",
			},
			expected: "[]()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hyperlink.String(); got != tt.expected {
				t.Errorf("Hyperlink.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
