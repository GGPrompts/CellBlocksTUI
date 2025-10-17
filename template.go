package main

import (
	"regexp"
	"strings"
)

// template.go - Template Variable Parsing
// Purpose: Detect and fill {{variable}} patterns in card content

// ExtractVariables finds all {{variable}} patterns in content
// Returns unique variable names (without braces or defaults)
func ExtractVariables(content string) []string {
	// Match {{variable}} or {{variable|default}}
	re := regexp.MustCompile(`\{\{([^}|]+)(?:\|[^}]*)?\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	// Collect unique variable names
	seen := make(map[string]bool)
	var vars []string

	for _, match := range matches {
		if len(match) > 1 {
			varName := strings.TrimSpace(match[1])
			if !seen[varName] {
				seen[varName] = true
				vars = append(vars, varName)
			}
		}
	}

	return vars
}

// ParseDefaultValue splits "variable|default" into (name, default)
// If no default is present, returns (name, "")
func ParseDefaultValue(variable string) (name, defaultVal string) {
	parts := strings.SplitN(variable, "|", 2)
	name = strings.TrimSpace(parts[0])

	if len(parts) > 1 {
		defaultVal = strings.TrimSpace(parts[1])
	}

	return name, defaultVal
}

// FillTemplate replaces all {{variable}} patterns with values from vars map
// If a variable has no value in the map, it uses the default value from {{var|default}}
// If no default and no value, keeps the original {{variable}}
func FillTemplate(content string, vars map[string]string) string {
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	result := re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract variable name (without braces)
		varWithDefault := strings.TrimSpace(match[2 : len(match)-2]) // Remove {{ and }}
		varName, defaultVal := ParseDefaultValue(varWithDefault)

		// Check if user provided a value
		if val, ok := vars[varName]; ok && val != "" {
			return val
		}

		// Fall back to default value
		if defaultVal != "" {
			return defaultVal
		}

		// Keep original if no value and no default
		return match
	})

	return result
}

// HasTemplateVariables checks if content contains any {{variable}} patterns
func HasTemplateVariables(content string) bool {
	re := regexp.MustCompile(`\{\{[^}]+\}\}`)
	return re.MatchString(content)
}
