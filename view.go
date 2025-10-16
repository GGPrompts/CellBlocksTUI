package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// view.go - UI Rendering
// Purpose: Render all UI components

// View is the main render function (Bubbletea lifecycle)
func (m Model) View() string {
	if m.Error != nil {
		return renderError(m.Error)
	}

	if m.ShowHelp {
		return renderHelp(m)
	}

	if m.Data == nil {
		return "Loading cards..."
	}

	var sections []string

	// Header
	sections = append(sections, renderHeader(m))

	// Main content area
	if m.ViewMode == ViewGrid {
		// Grid view
		sections = append(sections, renderGridView(m))
	} else if m.ShowPreview {
		// Split view: list + preview
		sections = append(sections, renderSplitView(m))
	} else {
		// Full list view
		sections = append(sections, renderListView(m))
	}

	// Status bar
	sections = append(sections, renderStatusBar(m))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderHeader renders the app title, search box, and card count
func renderHeader(m Model) string {
	title := styleTitle.Render("CellBlocks TUI")

	// Search query
	searchText := ""
	if m.SearchQuery != "" {
		searchText = styleSearchBox.Render(fmt.Sprintf(" Search: %s", m.SearchQuery))
	}

	// Card count
	count := styleSubtle.Render(fmt.Sprintf("[%d/%d]", len(m.FilteredCards), len(m.Data.Cards)))

	header := lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		searchText,
		" ",
		count,
	)

	return styleHeader.Width(m.Width).Render(header)
}

// renderListView renders the scrollable card list
func renderListView(m Model) string {
	return renderListViewWithHeight(m, m.Height-6)
}

// renderListViewWithHeight renders list with explicit height
func renderListViewWithHeight(m Model, availableHeight int) string {
	if len(m.FilteredCards) == 0 {
		return styleSubtle.Render("No cards found. Press '/' to clear filters or type to search.")
	}

	visibleCount := max(1, availableHeight)
	var lines []string

	// Calculate which cards to show
	start := m.ScrollOffset
	end := min(start+visibleCount, len(m.FilteredCards))

	for i := start; i < end; i++ {
		card := m.FilteredCards[i]
		isSelected := i == m.SelectedIndex

		lines = append(lines, renderCardListItem(m, &card, isSelected))
	}

	// Pad with empty lines if needed
	for len(lines) < visibleCount && len(lines) < len(m.FilteredCards) {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// renderCardListItem renders a single card in the list
func renderCardListItem(m Model, card *Card, selected bool) string {
	// Get category
	cat := m.getCategoryForCard(card)
	categoryName := ""
	categoryColor := ""
	if cat != nil {
		categoryName = cat.Name
		categoryColor = cat.Color
	}

	// Title
	title := card.Title
	maxTitleLen := m.Width - 20 // Leave space for category badge
	if len(title) > maxTitleLen {
		title = truncate(title, maxTitleLen)
	}

	// Category badge
	categoryBadge := styleCategoryName(categoryName, categoryColor)

	// Selection indicator
	indicator := " "
	if selected {
		indicator = styleCardTitleSelected.Render(">")
	}

	// Build line
	line := fmt.Sprintf("%s %s %s",
		indicator,
		title,
		categoryBadge,
	)

	if selected {
		return styleCardItemSelected.Width(m.Width).Render(line)
	}
	return styleCardItem.Width(m.Width).Render(line)
}

// renderSplitView renders list view on top, preview on bottom
func renderSplitView(m Model) string {
	availableHeight := m.Height - 6 // Header + status bar

	// Adaptive split: on taller screens (>50 lines), give more to preview
	// Mobile/small screens stay 50/50, desktop gets 40/60 list/preview
	var listHeight, previewHeight int
	if availableHeight > 50 {
		// Desktop: 40% list, 60% preview
		listHeight = availableHeight * 2 / 5
		previewHeight = availableHeight - listHeight
	} else {
		// Mobile: 50/50 split
		listHeight = availableHeight / 2
		previewHeight = availableHeight - listHeight
	}

	// Render list with explicit height
	listView := renderListViewWithHeight(m, listHeight)

	// Preview pane
	previewView := renderPreviewPane(m, previewHeight)

	return lipgloss.JoinVertical(lipgloss.Left, listView, previewView)
}

// renderGridView renders cards in a grid layout with neon borders
func renderGridView(m Model) string {
	if m.ShowPreview {
		return renderGridWithPreview(m)
	}

	if len(m.FilteredCards) == 0 {
		return styleSubtle.Render("No cards found. Press '/' to clear filters or type to search.")
	}

	return renderGridCards(m, m.Width, m.Height-6)
}

// renderGridCards renders grid cards within specified dimensions
func renderGridCards(m Model, width, height int) string {
	if len(m.FilteredCards) == 0 {
		return ""
	}

	// Calculate grid dimensions using constants
	cols := max(1, width/GridCardTotalWidth)
	rows := max(1, height/GridCardTotalHeight)

	// Calculate visible cards based on scroll offset
	cardsPerPage := cols * rows
	startIdx := (m.ScrollOffset / cols) * cols // Align to row boundary
	endIdx := min(startIdx+cardsPerPage, len(m.FilteredCards))

	// Build grid
	var gridRows []string
	currentRow := []string{}

	for i := startIdx; i < endIdx; i++ {
		card := m.FilteredCards[i]
		isSelected := i == m.SelectedIndex

		// Render card
		cardView := renderGridCard(m, &card, isSelected)
		currentRow = append(currentRow, cardView)

		// Complete row
		if len(currentRow) == cols {
			gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
			currentRow = []string{}
		}
	}

	// Add remaining cards in partial row
	if len(currentRow) > 0 {
		// Pad with empty cards if needed
		for len(currentRow) < cols {
			currentRow = append(currentRow, strings.Repeat(" ", GridCardTotalWidth))
		}
		gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, gridRows...)
}

// renderGridCard renders a single card in grid view
func renderGridCard(m Model, card *Card, selected bool) string {
	// Get category
	cat := m.getCategoryForCard(card)
	categoryColor := ""
	if cat != nil {
		categoryColor = cat.Color
	}

	// Wrap title to 3 lines (23 chars = 27 width - 4 for padding)
	lines := wrapText(card.Title, 23, 3)

	// Build card content
	content := strings.Join(lines, "\n")

	// Apply style with category-colored border
	cardStyle := makeGridCardStyle(categoryColor, selected)
	return cardStyle.Render(content)
}

// renderGridWithPreview renders grid with preview pane
func renderGridWithPreview(m Model) string {
	availableHeight := m.Height - 6 // Header + status bar

	// For wide screens (>120 chars), do side-by-side
	// For narrow screens, do top/bottom
	if m.Width > 120 {
		// Side-by-side: grid on left, preview on right
		gridWidth := m.Width * 3 / 5  // 60% for grid
		previewWidth := m.Width - gridWidth

		gridView := renderGridCards(m, gridWidth, availableHeight)
		previewView := renderPreviewPane(m, availableHeight)

		// Adjust preview width
		previewView = lipgloss.NewStyle().Width(previewWidth).Render(previewView)

		return lipgloss.JoinHorizontal(lipgloss.Top, gridView, previewView)
	} else {
		// Top/bottom: grid on top, preview on bottom
		var gridHeight, previewHeight int
		if availableHeight > 50 {
			gridHeight = availableHeight * 2 / 5      // 40% grid
			previewHeight = availableHeight - gridHeight // 60% preview
		} else {
			gridHeight = availableHeight / 2
			previewHeight = availableHeight - gridHeight
		}

		gridView := renderGridCards(m, m.Width, gridHeight)
		previewView := renderPreviewPane(m, previewHeight)

		return lipgloss.JoinVertical(lipgloss.Left, gridView, previewView)
	}
}

// renderPreviewPane renders the selected card's full content
func renderPreviewPane(m Model, height int) string {
	card := m.getSelectedCard()
	if card == nil {
		return stylePreviewPane.
			Width(m.Width - 2).
			Height(height - 2).
			Render("No card selected")
	}

	// Get category
	cat := m.getCategoryForCard(card)
	categoryName := ""
	categoryColor := ""
	if cat != nil {
		categoryName = cat.Name
		categoryColor = cat.Color
	}

	// Title with category
	title := stylePreviewTitle.Render(card.Title)
	category := styleCategoryName(categoryName, categoryColor)
	header := fmt.Sprintf("%s  %s", title, category)

	// Content - use almost all available space
	content := card.Content
	// More accurate calculation: border (2) + title (1) + blank line (1) + padding (2) = 6 total
	// But the Height() call already accounts for border, so we only need to subtract inner elements
	maxLines := height - 4 // Just title, blank line, and inner padding
	if maxLines < 3 {
		maxLines = 3 // Minimum 3 lines on tiny screens
	}

	contentLines := strings.Split(content, "\n")
	if len(contentLines) > maxLines {
		contentLines = contentLines[:maxLines]
		contentLines = append(contentLines, styleSubtle.Render("..."))
	}
	content = strings.Join(contentLines, "\n")

	preview := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		stylePreviewContent.Render(content),
	)

	return stylePreviewPane.
		Width(m.Width - 2).
		Height(height - 2).
		Render(preview)
}

// renderStatusBar renders keyboard shortcuts and help
func renderStatusBar(m Model) string {
	var hints []string

	if m.ViewMode == ViewGrid {
		hints = []string{
			styleHelpKey.Render("↑↓←→") + styleHelpDesc.Render(" navigate"),
			styleHelpKey.Render("Enter") + styleHelpDesc.Render(" copy"),
			styleHelpKey.Render("g") + styleHelpDesc.Render(" list"),
			styleHelpKey.Render("p") + styleHelpDesc.Render(" preview"),
			styleHelpKey.Render("?") + styleHelpDesc.Render(" help"),
			styleHelpKey.Render("q") + styleHelpDesc.Render(" quit"),
		}
	} else {
		hints = []string{
			styleHelpKey.Render("↑↓") + styleHelpDesc.Render(" navigate"),
			styleHelpKey.Render("Enter") + styleHelpDesc.Render(" copy"),
			styleHelpKey.Render("g") + styleHelpDesc.Render(" grid"),
			styleHelpKey.Render("p") + styleHelpDesc.Render(" preview"),
			styleHelpKey.Render("?") + styleHelpDesc.Render(" help"),
			styleHelpKey.Render("q") + styleHelpDesc.Render(" quit"),
		}
	}

	status := strings.Join(hints, "  ")
	return styleStatusBar.Width(m.Width).Render(" " + status)
}

// renderHelp renders the help dialog
func renderHelp(m Model) string {
	help := []string{
		styleTitle.Render("CellBlocks TUI - Keyboard Shortcuts"),
		"",
		styleHelpKey.Render("Navigation:"),
		"  ↑/k, ↓/j       Navigate cards",
		"  ←/h, →/l       Navigate left/right (grid view)",
		"  PgUp/PgDn      Scroll by page",
		"  Home/End       Jump to first/last",
		"",
		styleHelpKey.Render("View:"),
		"  g              Toggle grid/list view",
		"  p              Toggle preview pane (both modes)",
		"                 Side-by-side on wide screens!",
		"",
		styleHelpKey.Render("Actions:"),
		"  Enter, c       Copy card to clipboard",
		"  /              Clear filters",
		"  Type...        Search cards",
		"  Backspace      Delete search character",
		"",
		styleHelpKey.Render("General:"),
		"  ?              Toggle this help",
		"  Esc            Close help/clear search",
		"  q, Ctrl+C      Quit",
		"",
		styleSubtle.Render("Press ? or Esc to close"),
	}

	content := strings.Join(help, "\n")
	box := styleHelpBox.Render(content)

	// Center the box
	return lipgloss.Place(m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		box)
}

// renderError renders an error message
func renderError(err error) string {
	return styleError.Render(fmt.Sprintf("Error: %v\n\nPress q to quit.", err))
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
