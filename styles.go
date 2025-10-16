package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// styles.go - UI Styling with Lipgloss
// Purpose: Define all colors, borders, and text styles

var (
	// Colors
	colorPrimary   = lipgloss.Color("#00ff41") // Neon green
	colorSecondary = lipgloss.Color("#00ffff") // Cyan
	colorAccent    = lipgloss.Color("#ff00ff") // Magenta
	colorYellow    = lipgloss.Color("#ffff00") // Yellow
	colorOrange    = lipgloss.Color("#ff9500") // Orange
	colorBlue      = lipgloss.Color("#00a6ff") // Blue
	colorGray      = lipgloss.Color("#666666") // Gray
	colorDimmed    = lipgloss.Color("#444444") // Dimmed

	// Text styles
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	styleSubtle = lipgloss.NewStyle().
			Foreground(colorGray)

	styleDimmed = lipgloss.NewStyle().
			Foreground(colorDimmed)

	// Header
	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(colorPrimary)

	// Search box
	styleSearchBox = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	// Card list item
	styleCardItem = lipgloss.NewStyle().
			Padding(0, 1)

	styleCardItemSelected = lipgloss.NewStyle().
				Padding(0, 1).
				Background(lipgloss.Color("#1a1a1a")).
				Foreground(colorPrimary).
				Bold(true)

	// Card title in list
	styleCardTitle = lipgloss.NewStyle().
			Foreground(colorPrimary)

	styleCardTitleSelected = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true)

	// Category badge
	styleCategoryBadge = lipgloss.NewStyle().
				Padding(0, 1).
				MarginLeft(1).
				Background(lipgloss.Color("#222222")).
				Foreground(colorSecondary)

	// Preview pane
	stylePreviewPane = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(colorSecondary).
				Padding(1)

	stylePreviewTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorPrimary).
				Underline(true)

	stylePreviewContent = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff"))

	// Status bar
	styleStatusBar = lipgloss.NewStyle().
			Foreground(colorGray).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(colorGray)

	// Error message
	styleError = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Bold(true)

	// Help dialog
	styleHelpBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2).
			Background(lipgloss.Color("#0a0a0a"))

	styleHelpKey = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	styleHelpDesc = lipgloss.NewStyle().
			Foreground(colorGray)

	// Grid view card - compact 27 chars wide (was 36)
	styleGridCard = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1).
			Width(27).
			Height(4)

	styleGridCardSelected = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				Padding(0, 1).
				Width(27).
				Height(4).
				Background(lipgloss.Color("#1a1a1a")).
				Foreground(colorPrimary).
				Bold(true)
)

// getCategoryColor returns the color for a category
func getCategoryColor(color string) lipgloss.Color {
	if color == "" {
		return colorGray
	}
	return lipgloss.Color(color)
}

// styleCategoryName returns a styled category name with its color
func styleCategoryName(name string, color string) string {
	style := lipgloss.NewStyle().
		Foreground(getCategoryColor(color)).
		Bold(true)
	return style.Render(name)
}

// truncate truncates a string to maxLen and adds "..." if needed
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// wrapText wraps text to fit within maxWidth, returning up to maxLines
func wrapText(text string, maxWidth int, maxLines int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	currentLine := ""

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) > maxWidth {
			if currentLine == "" {
				// Word is too long, truncate it
				lines = append(lines, truncate(word, maxWidth))
				currentLine = ""
			} else {
				// Start new line with this word
				lines = append(lines, currentLine)
				currentLine = word
			}

			if len(lines) >= maxLines {
				break
			}
		} else {
			currentLine = testLine
		}
	}

	// Add remaining text
	if currentLine != "" && len(lines) < maxLines {
		lines = append(lines, currentLine)
	}

	// Pad to maxLines
	for len(lines) < maxLines {
		lines = append(lines, "")
	}

	return lines[:maxLines]
}

// makeGridCardStyle creates a card style with category-colored border
func makeGridCardStyle(categoryColor string, selected bool) lipgloss.Style {
	color := getCategoryColor(categoryColor)

	if selected {
		return styleGridCardSelected.
			BorderForeground(color)
	}
	return styleGridCard.
		BorderForeground(color)
}

// Grid card dimensions (keep these in sync with styleGridCard)
const (
	GridCardWidth  = 27 // Inner width
	GridCardHeight = 4  // Inner height
	GridCardSpacing = 2 // Border + spacing
	GridCardTotalWidth = GridCardWidth + GridCardSpacing  // 29 total
	GridCardTotalHeight = GridCardHeight + GridCardSpacing // 6 total
)
