package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
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

	// Category filter screen
	if m.ViewMode == ViewCategoryFilter {
		return renderCategoryFilterScreen(m)
	}

	// Card creation screen
	if m.ViewMode == ViewCardCreate {
		return renderCardCreateScreen(m)
	}

	// Detail view (full-screen card)
	if m.ViewMode == ViewDetail {
		return renderDetailView(m)
	}

	var sections []string

	// Header
	sections = append(sections, renderHeader(m))

	// Main content area
	if m.ViewMode == ViewGrid {
		// Grid view
		sections = append(sections, renderGridView(m))
	} else if m.ViewMode == ViewTable {
		// Table view
		sections = append(sections, renderTableView(m))
	} else if m.ShowPreview {
		// Split view: list + preview
		sections = append(sections, renderSplitView(m))
	} else {
		// Full list view
		sections = append(sections, renderListView(m))
	}

	// Status bar
	sections = append(sections, renderStatusBar(m))

	// Join all sections and ensure they fill the screen height to prevent ghosting
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Use Place to ensure content fills entire screen and clears any previous render artifacts
	return lipgloss.Place(m.Width, m.Height,
		lipgloss.Left, lipgloss.Top,
		content,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.NoColor{}),
	)
}

// renderHeader renders the app title, search box, and card count
func renderHeader(m Model) string {
	var headerLines []string

	// Main header line
	title := styleTitle.Render("CellBlocks TUI")

	// Category filters - mobile-friendly display
	filterText := ""
	if len(m.SelectedCategories) > 0 {
		// On narrow screens (< 80 chars), just show count
		// On wider screens, show first 2 categories + count
		if m.Width < 80 {
			filterText = styleSearchBox.Render(fmt.Sprintf(" Filters: %d", len(m.SelectedCategories)))
		} else {
			// Collect all category names and sort them for stable display
			var catNames []string
			for catID := range m.SelectedCategories {
				if cat, ok := m.CategoryMap[catID]; ok {
					catNames = append(catNames, cat.Name)
				}
			}
			// Sort alphabetically to prevent flickering
			sortStrings(catNames)

			// Show first 2 category names
			var displayNames []string
			for i := 0; i < len(catNames) && i < 2; i++ {
				displayNames = append(displayNames, catNames[i])
			}

			remaining := len(catNames) - len(displayNames)
			if remaining > 0 {
				filterText = styleSearchBox.Render(fmt.Sprintf(" Filters: %s +%d", strings.Join(displayNames, ", "), remaining))
			} else {
				filterText = styleSearchBox.Render(fmt.Sprintf(" Filters: %s", strings.Join(displayNames, ", ")))
			}
		}
	}

	// Card count
	count := styleSubtle.Render(fmt.Sprintf("[%d/%d]", len(m.FilteredCards), len(m.Data.Cards)))

	mainLine := lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		filterText,
		" ",
		count,
	)
	headerLines = append(headerLines, mainLine)

	// Reload notification (show for 5 seconds)
	if m.ReloadMessage != "" && time.Since(m.ReloadMessageTime) < 5*time.Second {
		notifStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff41")).
			Bold(true)
		headerLines = append(headerLines, notifStyle.Render(m.ReloadMessage))
	}

	header := lipgloss.JoinVertical(lipgloss.Left, headerLines...)
	return styleHeader.Width(m.Width).Render(header)
}

// renderListView renders the scrollable card list
func renderListView(m Model) string {
	return renderListViewWithHeight(m, m.Height-6)
}

// renderListViewWithHeight renders list with explicit height
func renderListViewWithHeight(m Model, availableHeight int) string {
	if len(m.FilteredCards) == 0 {
		return styleSubtle.Render("No cards found. Press 'f' to filter by category.")
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
	if runewidth.StringWidth(title) > maxTitleLen {
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
		return styleSubtle.Render("No cards found. Press 'f' to filter by category.")
	}

	return renderGridCards(m, m.Width, m.Height-6)
}

// renderGridCards renders grid cards within specified dimensions
func renderGridCards(m Model, width, height int) string {
	if len(m.FilteredCards) == 0 {
		return ""
	}

	// Validate dimensions to prevent panics
	if width < GridCardTotalWidth {
		width = GridCardTotalWidth
	}
	if height < GridCardTotalHeight {
		height = GridCardTotalHeight
	}

	// Calculate grid dimensions using constants
	cols := max(1, min(width/GridCardTotalWidth, GridMaxColumns))
	rows := max(1, height/GridCardTotalHeight)

	// Calculate visible cards based on scroll offset
	cardsPerPage := cols * rows
	startIdx := (m.ScrollOffset / cols) * cols // Align to row boundary

	// Sanity check: ensure startIdx is valid
	if startIdx < 0 {
		startIdx = 0
	}
	if startIdx >= len(m.FilteredCards) {
		startIdx = 0
	}

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

	// Add remaining cards in partial row (without padding empty spaces)
	if len(currentRow) > 0 {
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

	// Card has 6 lines available (GridCardHeight = 6)
	// Strategy: Title (1-2 lines) + Content preview (remaining lines)

	// Wrap title to max 2 lines (25 chars = 27 width - 2 for horizontal padding)
	titleLines := wrapText(card.Title, 25, 2)

	// Calculate remaining lines for content preview
	remainingLines := 6 - len(titleLines)

	// Build output
	var lines []string

	// Add title (bold/primary)
	for _, line := range titleLines {
		if selected {
			lines = append(lines, styleCardTitleSelected.Render(line))
		} else {
			lines = append(lines, styleCardTitle.Render(line))
		}
	}

	// Add content preview (dimmed) if we have space
	if remainingLines > 0 && card.Content != "" {
		// Clean content: strip newlines, truncate
		contentPreview := strings.ReplaceAll(card.Content, "\n", " ")
		contentPreview = strings.TrimSpace(contentPreview)

		// Wrap content to remaining lines (25 chars to match title width)
		contentLines := wrapText(contentPreview, 25, remainingLines)
		for _, line := range contentLines {
			lines = append(lines, styleSubtle.Render(line))
		}
	}

	// Pad to 6 lines if needed
	for len(lines) < 6 {
		lines = append(lines, "")
	}

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
		// Side-by-side: FIXED 50/50 split for stability
		// This avoids complex calculations that can produce invalid dimensions
		gridWidth := m.Width / 2
		previewWidth := m.Width - gridWidth - 4 // Account for spacing and borders

		// Ensure minimum widths
		if gridWidth < GridCardTotalWidth {
			gridWidth = GridCardTotalWidth
		}
		if previewWidth < 30 {
			previewWidth = 30
		}

		gridView := renderGridCards(m, gridWidth, availableHeight)
		previewView := renderPreviewPaneWithWidth(m, availableHeight, previewWidth)

		return lipgloss.JoinHorizontal(lipgloss.Top, gridView, "  ", previewView)
	} else {
		// Top/bottom: FIXED 40/60 split for stability
		gridHeight := (availableHeight * 2) / 5      // 40% for grid
		previewHeight := availableHeight - gridHeight // 60% for preview

		// Ensure minimum heights
		if gridHeight < GridCardTotalHeight {
			gridHeight = GridCardTotalHeight
		}
		if previewHeight < 10 {
			previewHeight = 10
		}

		gridView := renderGridCards(m, m.Width, gridHeight)

		// Ensure grid fills allocated space to prevent preview showing through
		gridStyle := lipgloss.NewStyle().
			Height(gridHeight).
			Width(m.Width)
		gridView = gridStyle.Render(gridView)

		previewView := renderPreviewPane(m, previewHeight)

		return lipgloss.JoinVertical(lipgloss.Left, gridView, previewView)
	}
}

// renderTableView renders cards in an Excel-style table with sortable columns
func renderTableView(m Model) string {
	if len(m.FilteredCards) == 0 {
		return styleSubtle.Render("No cards found. Press 'f' to filter by category.")
	}

	// Sort cards based on current sort column and direction
	sortedCards := sortCards(m.FilteredCards, m.CategoryMap, m.SortColumn, m.SortDirection)

	// Calculate column widths based on terminal width
	// Available width = terminal width - borders and padding
	availableWidth := m.Width - 4

	// Column width distribution (percentages of available width)
	// Title: 40%, Category: 20%, Created: 20%, Updated: 20%
	titleWidth := availableWidth * 4 / 10
	categoryWidth := availableWidth * 2 / 10
	createdWidth := availableWidth * 2 / 10
	updatedWidth := availableWidth - titleWidth - categoryWidth - createdWidth // Remaining space

	// Minimum widths to prevent squishing
	if titleWidth < 20 {
		titleWidth = 20
	}
	if categoryWidth < 10 {
		categoryWidth = 10
	}
	if createdWidth < 10 {
		createdWidth = 10
	}
	if updatedWidth < 10 {
		updatedWidth = 10
	}

	// Build header row with sort indicators
	titleHeader := "Title" + getSortIndicator("title", m.SortColumn, m.SortDirection)
	categoryHeader := "Category" + getSortIndicator("category", m.SortColumn, m.SortDirection)
	createdHeader := "Created" + getSortIndicator("created", m.SortColumn, m.SortDirection)
	updatedHeader := "Updated" + getSortIndicator("updated", m.SortColumn, m.SortDirection)

	// Pad headers to column width
	titleHeader = padOrTruncate(titleHeader, titleWidth)
	categoryHeader = padOrTruncate(categoryHeader, categoryWidth)
	createdHeader = padOrTruncate(createdHeader, createdWidth)
	updatedHeader = padOrTruncate(updatedHeader, updatedWidth)

	// Style the header (with 2-space indent to match data rows)
	headerRow := styleTableHeader.Render(
		fmt.Sprintf("  %s │ %s │ %s │ %s",
			titleHeader,
			categoryHeader,
			createdHeader,
			updatedHeader,
		),
	)

	// Separator line (with 2-space indent to match data rows)
	separator := "  " + strings.Repeat("─", titleWidth) + "─┼─" +
		strings.Repeat("─", categoryWidth) + "─┼─" +
		strings.Repeat("─", createdWidth) + "─┼─" +
		strings.Repeat("─", updatedWidth)

	var lines []string
	lines = append(lines, headerRow)
	lines = append(lines, separator)

	// Calculate visible rows
	availableHeight := m.Height - 6 - 2 // Header + status bar + table header
	visibleCount := max(1, availableHeight)

	// Calculate which cards to show based on scroll offset
	start := m.ScrollOffset
	end := min(start+visibleCount, len(sortedCards))

	// Render table rows
	for i := start; i < end; i++ {
		card := sortedCards[i]
		isSelected := i == m.SelectedIndex

		// Get category name
		categoryName := ""
		categoryColor := ""
		if cat, ok := m.CategoryMap[card.CategoryID]; ok {
			categoryName = cat.Name
			categoryColor = cat.Color
		}

		// Format data for display
		title := padOrTruncate(card.Title, titleWidth)
		category := padOrTruncate(categoryName, categoryWidth)
		created := padOrTruncate(formatDate(card.CreatedAt), createdWidth)
		updated := padOrTruncate(formatDate(card.UpdatedAt), updatedWidth)

		// Build row
		row := fmt.Sprintf("%s │ %s │ %s │ %s",
			title,
			styleCategoryName(category, categoryColor),
			created,
			updated,
		)

		// Apply selection style
		if isSelected {
			row = styleCardItemSelected.Width(m.Width).Render("> " + row)
		} else {
			row = styleCardItem.Width(m.Width).Render("  " + row)
		}

		lines = append(lines, row)
	}

	// Pad with empty lines if needed
	for len(lines)-2 < visibleCount { // -2 for header and separator
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// padOrTruncate pads a string to the specified width or truncates if too long
// Uses display width (not byte length) to handle emojis and wide characters correctly
func padOrTruncate(s string, width int) string {
	displayWidth := runewidth.StringWidth(s)

	if displayWidth > width {
		if width <= 3 {
			return runewidth.Truncate(s, width, "")
		}
		return runewidth.Truncate(s, width-3, "") + "..."
	}

	// Pad with spaces to reach the target width
	padding := width - displayWidth
	if padding > 0 {
		return s + strings.Repeat(" ", padding)
	}
	return s
}

// renderPreviewPane renders the selected card's full content
func renderPreviewPane(m Model, height int) string {
	return renderPreviewPaneWithWidth(m, height, m.Width-2)
}

// renderPreviewPaneWithWidth renders preview pane with custom width
func renderPreviewPaneWithWidth(m Model, height int, width int) string {
	// Validate dimensions to prevent panics from glamour/lipgloss
	if width < 10 {
		width = 10
	}
	if height < 5 {
		height = 5
	}

	card := m.getPreviewedCard()
	if card == nil {
		return stylePreviewPane.
			Width(width).
			Height(height - 2).
			Render("No card previewed - click a card to preview")
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

	// Render with markdown if enabled
	if m.UseMarkdownRender {
		content = renderMarkdown(content, width-4)
	}

	// More accurate calculation: border (2) + title (1) + blank line (1) + padding (2) = 6 total
	// But the Height() call already accounts for border, so we only need to subtract inner elements
	maxLines := height - 4 // Just title, blank line, and inner padding
	if maxLines < 3 {
		maxLines = 3 // Minimum 3 lines on tiny screens
	}

	contentLines := strings.Split(content, "\n")
	totalLines := len(contentLines)

	// Apply scroll offset
	startLine := m.PreviewScrollOffset
	if startLine >= totalLines {
		startLine = max(0, totalLines-1)
	}

	endLine := min(startLine+maxLines, totalLines)
	visibleLines := contentLines[startLine:endLine]

	// Add scroll indicators
	scrollInfo := ""
	if startLine > 0 {
		scrollInfo = styleSubtle.Render(fmt.Sprintf("▲ (line %d/%d)", startLine+1, totalLines))
	}
	if endLine < totalLines {
		if scrollInfo != "" {
			scrollInfo += " "
		}
		scrollInfo += styleSubtle.Render(fmt.Sprintf("▼ (Shift+↑↓ to scroll)"))
	}

	if scrollInfo != "" {
		visibleLines = append(visibleLines, scrollInfo)
	}

	content = strings.Join(visibleLines, "\n")

	preview := lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		stylePreviewContent.Render(content),
	)

	return stylePreviewPane.
		Width(width).
		Height(height - 2).
		Render(preview)
}

// renderStatusBar renders keyboard shortcuts and help
func renderStatusBar(m Model) string {
	var hints []string

	// Responsive hints based on screen width
	if m.Width < 60 {
		// Very narrow (mobile) - show only essentials
		hints = []string{
			styleHelpKey.Render("n") + styleHelpDesc.Render(" new"),
			styleHelpKey.Render("f") + styleHelpDesc.Render(" filter"),
			styleHelpKey.Render("?") + styleHelpDesc.Render(" help"),
			styleHelpKey.Render("q") + styleHelpDesc.Render(" quit"),
		}
	} else if m.Width < 90 {
		// Narrow - show most important
		hints = []string{
			styleHelpKey.Render("Enter") + styleHelpDesc.Render(" copy"),
			styleHelpKey.Render("n") + styleHelpDesc.Render(" new"),
			styleHelpKey.Render("f") + styleHelpDesc.Render(" filter"),
			styleHelpKey.Render("?") + styleHelpDesc.Render(" help"),
			styleHelpKey.Render("q") + styleHelpDesc.Render(" quit"),
		}
	} else {
		// Wide screen - show all hints
		if m.ViewMode == ViewGrid {
			hints = []string{
				styleHelpKey.Render("↑↓←→") + styleHelpDesc.Render(" navigate"),
				styleHelpKey.Render("Enter") + styleHelpDesc.Render(" copy"),
				styleHelpKey.Render("n") + styleHelpDesc.Render(" new"),
				styleHelpKey.Render("f") + styleHelpDesc.Render(" filter"),
				styleHelpKey.Render("g") + styleHelpDesc.Render(" table"),
				styleHelpKey.Render("?") + styleHelpDesc.Render(" help"),
			}
		} else if m.ViewMode == ViewTable {
			hints = []string{
				styleHelpKey.Render("↑↓") + styleHelpDesc.Render(" navigate"),
				styleHelpKey.Render("1-4") + styleHelpDesc.Render(" sort"),
				styleHelpKey.Render("Enter") + styleHelpDesc.Render(" view"),
				styleHelpKey.Render("n") + styleHelpDesc.Render(" new"),
				styleHelpKey.Render("g") + styleHelpDesc.Render(" list"),
				styleHelpKey.Render("?") + styleHelpDesc.Render(" help"),
			}
		} else {
			hints = []string{
				styleHelpKey.Render("↑↓") + styleHelpDesc.Render(" navigate"),
				styleHelpKey.Render("Enter") + styleHelpDesc.Render(" copy"),
				styleHelpKey.Render("n") + styleHelpDesc.Render(" new"),
				styleHelpKey.Render("f") + styleHelpDesc.Render(" filter"),
				styleHelpKey.Render("g") + styleHelpDesc.Render(" grid"),
				styleHelpKey.Render("?") + styleHelpDesc.Render(" help"),
			}
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
		"  Shift+↑/↓      Scroll preview content",
		"  PgUp/PgDn      Scroll by page",
		"  Home/End       Jump to first/last",
		"",
		styleHelpKey.Render("View:"),
		"  g              Cycle through list → grid → table view",
		"  p              Toggle preview pane (list/grid modes)",
		"                 Side-by-side on wide screens!",
		"  Space          Pin card to preview (grid view)",
		"",
		styleHelpKey.Render("Table View:"),
		"  1              Sort by title (press again to reverse)",
		"  2              Sort by category",
		"  3              Sort by created date",
		"  4              Sort by updated date",
		"  ↑↓/k/j         Navigate rows",
		"",
		styleHelpKey.Render("Actions:"),
		"  Enter, d       Open card in detail view",
		"  c              Copy card to clipboard",
		"  n              Create new card",
		"  f              Filter by category",
		"",
		styleHelpKey.Render("Detail View:"),
		"  ↑/↓, k/j       Scroll content",
		"  m              Toggle markdown rendering",
		"  t              Toggle template form (if templates detected)",
		"  Tab            Navigate template fields",
		"  Enter, c       Copy (filled template if editing)",
		"  Esc            Return to list/grid view",
		"",
		styleHelpKey.Render("Mouse/Touch:"),
		"  Click          Select & pin to preview (grid)",
		"  Double-click   Copy card to clipboard",
		"  Mouse wheel    Scroll preview (over preview pane)",
		"",
		styleHelpKey.Render("Auto-Reload:"),
		"  ✨             Checks for new cards every 10 seconds",
		"                 (Perfect for AI-generated cards!)",
		"",
		styleHelpKey.Render("General:"),
		"  ?              Toggle this help",
		"  Esc            Close help",
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

// sortStrings sorts a slice of strings in place (simple bubble sort for small lists)
func sortStrings(s []string) {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// renderCategoryFilterScreen renders the category filter selection screen
func renderCategoryFilterScreen(m Model) string {
	if m.Data == nil {
		return "No data loaded"
	}

	var lines []string

	// Title
	title := styleTitle.Render("Filter by Category")
	lines = append(lines, title)
	lines = append(lines, "")

	// Instructions
	instructions := styleSubtle.Render("↑↓: Navigate  Space/Enter: Toggle  A: All  C: Clear  Esc: Back")
	lines = append(lines, instructions)
	lines = append(lines, "")

	// Category list
	activeCount := len(m.SelectedCategories)
	filterInfo := fmt.Sprintf("Active filters: %d", activeCount)
	if activeCount > 0 {
		filterInfo = styleSearchBox.Render(filterInfo)
	} else {
		filterInfo = styleSubtle.Render(filterInfo)
	}
	lines = append(lines, filterInfo)
	lines = append(lines, "")

	// Render each category
	for i, cat := range m.Data.Categories {
		isSelected := m.FilterCursorIndex == i
		isActive := m.SelectedCategories[cat.ID]

		// Checkbox
		checkbox := "[ ]"
		if isActive {
			checkbox = "[✓]"
		}

		// Category name with color
		catName := styleCategoryName(cat.Name, cat.Color)

		// Build line
		var line string
		if isSelected {
			indicator := styleCardTitleSelected.Render(">")
			line = fmt.Sprintf("%s %s %s", indicator, checkbox, catName)
			line = styleCardItemSelected.Render(line)
		} else {
			line = fmt.Sprintf("  %s %s", checkbox, catName)
			line = styleCardItem.Render(line)
		}

		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")

	// Center on screen
	return lipgloss.Place(m.Width, m.Height,
		lipgloss.Center, lipgloss.Top,
		content)
}

// renderCardCreateScreen renders the card creation form
func renderCardCreateScreen(m Model) string {
	if m.Data == nil {
		return "No data loaded"
	}

	var lines []string

	// Title
	title := styleTitle.Render("Create New Card")
	lines = append(lines, title)
	lines = append(lines, "")

	// Instructions
	instructions := styleSubtle.Render("Tab: Next field  Ctrl+S: Save  Esc: Cancel")
	lines = append(lines, instructions)
	lines = append(lines, "")

	// Field 0: Title
	titleLabel := "Title:"
	if m.CreateFormField == 0 {
		titleLabel = styleSearchBox.Render("→ Title:")
	}
	lines = append(lines, titleLabel)
	titleValue := m.NewCardTitle
	if titleValue == "" {
		titleValue = styleSubtle.Render("(enter title)")
	}
	if m.CreateFormField == 0 {
		titleValue = styleCardItemSelected.Render(titleValue + "█") // Show cursor
	}
	lines = append(lines, "  "+titleValue)
	lines = append(lines, "")

	// Field 1: Content
	contentLabel := "Content:"
	if m.CreateFormField == 1 {
		contentLabel = styleSearchBox.Render("→ Content:")
	}
	lines = append(lines, contentLabel)

	// Show content (multi-line)
	contentValue := m.NewCardContent
	if contentValue == "" {
		contentValue = styleSubtle.Render("(enter content)")
		lines = append(lines, "  "+contentValue)
	} else {
		// Split by newlines and render
		contentLines := strings.Split(contentValue, "\n")
		maxLines := 10 // Limit visible lines
		for i, line := range contentLines {
			if i >= maxLines {
				lines = append(lines, "  "+styleSubtle.Render(fmt.Sprintf("... (%d more lines)", len(contentLines)-maxLines)))
				break
			}
			// Show cursor on last line if focused
			if m.CreateFormField == 1 && i == len(contentLines)-1 {
				line = styleCardItemSelected.Render(line + "█")
			}
			lines = append(lines, "  "+line)
		}
	}
	lines = append(lines, "")

	// Field 2: Category
	categoryLabel := "Category:"
	if m.CreateFormField == 2 {
		categoryLabel = styleSearchBox.Render("→ Category:")
	}
	lines = append(lines, categoryLabel)

	// Show current category selection
	var selectedCat *Category
	for _, cat := range m.Data.Categories {
		if cat.ID == m.NewCardCategoryID {
			selectedCat = &cat
			break
		}
	}
	if selectedCat != nil {
		catDisplay := styleCategoryName(selectedCat.Name, selectedCat.Color)
		if m.CreateFormField == 2 {
			catDisplay = styleCardItemSelected.Render(catDisplay + " (↑↓ to change)")
		}
		lines = append(lines, "  "+catDisplay)
	} else {
		lines = append(lines, "  "+styleSubtle.Render("(no category selected)"))
	}

	lines = append(lines, "")
	lines = append(lines, "")

	// Validation hints
	if m.NewCardTitle == "" || m.NewCardContent == "" {
		hint := styleError.Render("⚠ Title and content are required")
		lines = append(lines, hint)
	} else {
		hint := styleHelpKey.Render("✓ Ready to save! Press Ctrl+S")
		lines = append(lines, hint)
	}

	content := strings.Join(lines, "\n")

	// Center on screen
	return lipgloss.Place(m.Width, m.Height,
		lipgloss.Center, lipgloss.Top,
		content)
}

// renderDetailView renders full-screen card view with markdown and templates
func renderDetailView(m Model) string {
	card := m.getSelectedCard()
	if card == nil {
		return styleSubtle.Render("No card selected. Press Esc to return.")
	}

	// Get category
	cat := m.getCategoryForCard(card)
	categoryName := ""
	categoryColor := ""
	if cat != nil {
		categoryName = cat.Name
		categoryColor = cat.Color
	}

	var lines []string

	// Header: Title + Category + Markdown indicator + separator
	title := stylePreviewTitle.Render(card.Title)
	category := styleCategoryName(categoryName, categoryColor)

	var mdIndicator string
	if m.UseMarkdownRender {
		mdIndicator = styleSearchBox.Render(" [MD] ")
	} else {
		mdIndicator = styleSubtle.Render(" [TXT] ")
	}

	header := lipgloss.JoinHorizontal(lipgloss.Left, title, "  ", category, mdIndicator)

	separator := strings.Repeat("─", m.Width-4)

	lines = append(lines, header)
	lines = append(lines, separator)
	lines = append(lines, "")

	// Available height for content and template form
	// Header (title + separator + blank = 3) + footer (blank + footer = 2) = 5 lines total
	availableHeight := m.Height - 5

	// Render content with optional markdown
	content := card.Content
	var renderedContent string

	// Use most of screen width (leave margin)
	contentWidth := m.Width - 8

	if m.UseMarkdownRender {
		// Render markdown using glamour
		renderedContent = renderMarkdown(content, contentWidth)
	} else {
		renderedContent = content
	}

	// Check for template variables
	hasTemplates := HasTemplateVariables(content)

	// Split available space between content and template form
	var contentHeight int
	var templateHeight int

	if hasTemplates && m.ShowTemplateForm {
		// Calculate heights based on number of variables
		numVars := len(m.DetectedVars)
		// Template form needs: header (1) + blank (1) + vars (n*2) + blank (1) + preview header (1) + preview (3) = 7 + n*2
		templateFormLines := 7 + numVars*2
		templateHeight = min(templateFormLines, availableHeight/2)
		contentHeight = availableHeight - templateHeight
	} else {
		contentHeight = availableHeight
		templateHeight = 0
	}

	// Render scrollable content
	contentLines := strings.Split(renderedContent, "\n")
	totalLines := len(contentLines)

	maxContentLines := max(3, contentHeight)

	// Calculate maximum scroll offset (don't scroll past the end)
	maxScroll := max(0, totalLines-maxContentLines)

	// Clamp scroll offset to valid range
	startLine := m.DetailScrollOffset
	if startLine > maxScroll {
		startLine = maxScroll
	}
	if startLine < 0 {
		startLine = 0
	}

	endLine := min(startLine+maxContentLines, totalLines)

	// Handle case where content fits entirely on screen
	if totalLines <= maxContentLines {
		startLine = 0
		endLine = totalLines
	}

	visibleLines := contentLines[startLine:endLine]

	// Add scroll indicators
	scrollInfo := ""
	if startLine > 0 {
		scrollInfo = styleSubtle.Render(fmt.Sprintf("▲ Line %d/%d", startLine+1, totalLines))
	}
	if endLine < totalLines {
		if scrollInfo != "" {
			scrollInfo += " "
		}
		scrollInfo += styleSubtle.Render("▼ Scroll with ↑↓")
	}

	if scrollInfo != "" {
		visibleLines = append(visibleLines, "", scrollInfo)
	}

	lines = append(lines, visibleLines...)

	// Template form (if applicable)
	if hasTemplates && m.ShowTemplateForm {
		lines = append(lines, "")
		lines = append(lines, renderTemplateForm(m, card))
	}

	// Footer/instructions
	footer := buildDetailFooter(m, hasTemplates)
	lines = append(lines, "", footer)

	finalContent := strings.Join(lines, "\n")

	// Use full screen like other full-screen views
	return lipgloss.Place(m.Width, m.Height,
		lipgloss.Left, lipgloss.Top,
		finalContent)
}

// renderTemplateForm renders the template variable input form
func renderTemplateForm(m Model, card *Card) string {
	var lines []string

	// Title
	lines = append(lines, styleHelpKey.Render("Template Variables:"))
	lines = append(lines, "")

	// Render input fields for each variable
	for i, varName := range m.DetectedVars {
		isSelected := m.TemplateFormField == i

		// Label
		label := varName + ":"
		if isSelected {
			label = styleSearchBox.Render("→ " + label)
		} else {
			label = "  " + label
		}

		// Value
		value := m.TemplateVars[varName]
		if value == "" {
			value = styleSubtle.Render("(enter value)")
		}
		if isSelected {
			value = styleCardItemSelected.Render(value + "█")
		}

		lines = append(lines, label)
		lines = append(lines, "  "+value)
	}

	// Preview of filled template
	lines = append(lines, "")
	lines = append(lines, styleHelpKey.Render("Preview:"))
	filledContent := FillTemplate(card.Content, m.TemplateVars)

	// Show first few lines of filled content
	previewLines := strings.Split(filledContent, "\n")
	maxPreview := min(3, len(previewLines))
	for i := 0; i < maxPreview; i++ {
		lines = append(lines, stylePreviewContent.Render(previewLines[i]))
	}
	if len(previewLines) > maxPreview {
		lines = append(lines, styleSubtle.Render("..."))
	}

	return strings.Join(lines, "\n")
}

// buildDetailFooter creates the footer with keyboard shortcuts
func buildDetailFooter(m Model, hasTemplates bool) string {
	var hints []string

	if hasTemplates && m.ShowTemplateForm {
		hints = []string{
			styleHelpKey.Render("Tab") + styleHelpDesc.Render(" next field"),
			styleHelpKey.Render("Enter") + styleHelpDesc.Render(" copy filled"),
			styleHelpKey.Render("t") + styleHelpDesc.Render(" hide form"),
		}
	} else if hasTemplates {
		hints = []string{
			styleHelpKey.Render("t") + styleHelpDesc.Render(" show template form"),
			styleHelpKey.Render("c") + styleHelpDesc.Render(" copy"),
		}
	} else {
		hints = []string{
			styleHelpKey.Render("c") + styleHelpDesc.Render(" copy"),
		}
	}

	// Common shortcuts
	if m.UseMarkdownRender {
		hints = append(hints, styleHelpKey.Render("m") + styleHelpDesc.Render(" plain text"))
	} else {
		hints = append(hints, styleHelpKey.Render("m") + styleHelpDesc.Render(" markdown"))
	}

	hints = append(hints, styleHelpKey.Render("Esc") + styleHelpDesc.Render(" back"))

	return strings.Join(hints, "  ")
}

// renderMarkdown uses glamour to render markdown content
func renderMarkdown(content string, width int) string {
	// Validate width to prevent glamour panics
	if width < 10 {
		width = 10
	}

	// Create a new glamour renderer with dark style
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		// Fallback to plain text if glamour fails
		return content
	}

	// Render the markdown
	rendered, err := r.Render(content)
	if err != nil {
		// Fallback to plain text on error
		return content
	}

	// Trim trailing newlines that glamour adds
	return strings.TrimRight(rendered, "\n")
}
