package main

import (
	"reflect"
	"testing"
)

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "simple variable",
			content:  "docker run -p {{port}}:{{port}} {{image}}",
			expected: []string{"port", "image"},
		},
		{
			name:     "variable with default",
			content:  "server --port={{port|3000}} --host={{host|localhost}}",
			expected: []string{"port", "host"},
		},
		{
			name:     "no variables",
			content:  "plain text content",
			expected: []string{},
		},
		{
			name:     "duplicate variables",
			content:  "{{name}} says hello to {{name}}",
			expected: []string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractVariables(tt.content)
			// Handle nil vs empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExtractVariables(%q) = %v, want %v", tt.content, result, tt.expected)
			}
		})
	}
}

func TestParseDefaultValue(t *testing.T) {
	tests := []struct {
		input         string
		expectedName  string
		expectedValue string
	}{
		{"port|3000", "port", "3000"},
		{"host|localhost", "host", "localhost"},
		{"simple", "simple", ""},
		{"  spaced | value  ", "spaced", "value"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			name, value := ParseDefaultValue(tt.input)
			if name != tt.expectedName || value != tt.expectedValue {
				t.Errorf("ParseDefaultValue(%q) = (%q, %q), want (%q, %q)",
					tt.input, name, value, tt.expectedName, tt.expectedValue)
			}
		})
	}
}

func TestFillTemplate(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		vars     map[string]string
		expected string
	}{
		{
			name:    "fill all variables",
			content: "docker run -p {{port}}:{{port}} {{image}}",
			vars: map[string]string{
				"port":  "8080",
				"image": "nginx",
			},
			expected: "docker run -p 8080:8080 nginx",
		},
		{
			name:    "use default values",
			content: "server --port={{port|3000}} --host={{host|localhost}}",
			vars:    map[string]string{},
			expected: "server --port=3000 --host=localhost",
		},
		{
			name:    "partial fill with defaults",
			content: "server --port={{port|3000}} --host={{host|localhost}}",
			vars: map[string]string{
				"port": "8080",
			},
			expected: "server --port=8080 --host=localhost",
		},
		{
			name:     "no variables",
			content:  "plain text",
			vars:     map[string]string{},
			expected: "plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FillTemplate(tt.content, tt.vars)
			if result != tt.expected {
				t.Errorf("FillTemplate() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHasTemplateVariables(t *testing.T) {
	tests := []struct {
		content  string
		expected bool
	}{
		{"docker run {{image}}", true},
		{"plain text", false},
		{"{{var1}} and {{var2}}", true},
		{"{single brace}", false},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			result := HasTemplateVariables(tt.content)
			if result != tt.expected {
				t.Errorf("HasTemplateVariables(%q) = %v, want %v", tt.content, result, tt.expected)
			}
		})
	}
}
