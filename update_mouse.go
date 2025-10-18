package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// update_mouse.go - Enhanced Mouse/Touch Navigation
// Purpose: Handle mouse events with click, double-click, and scroll gestures
// Ported from TFE for CellBlocksTUI

// handleMouseEventEnhanced processes all mouse input with click and gesture detection
func (m Model) handleMouseEventEnhanced(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Handle ViewDetail mode separately (support scrolling content)
	if m.ViewMode == ViewDetail {
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.DetailScrollOffset = max(0, m.DetailScrollOffset-3)
			return m, nil
		case tea.MouseButtonWheelDown:
			// Only scroll if there's more content below
			// This will be bounds-checked in the render function
			m.DetailScrollOffset += 3
			return m, nil
		}
		// Ignore other mouse events in detail view
		return m, nil
	}

	// Don't process mouse events in filter/create screens
	if m.ViewMode == ViewCategoryFilter || m.ViewMode == ViewCardCreate {
		return m, nil
	}

	// Handle mouse wheel scrolling
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		// Check if mouse is over preview pane
		if m.ShowPreview && m.isMouseOverPreview(msg) {
			// Mouse is over preview - scroll preview content
			m.PreviewScrollOffset = max(0, m.PreviewScrollOffset-3)
			return m, nil
		}
		// Otherwise scroll card list
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, -1) // Move up one row
		} else {
			m.moveSelection(-1)
			// In list view with preview: keep preview locked (don't update PreviewedIndex)
			// Preview only updates on keyboard nav, click, or spacebar pin
		}
		return m, nil

	case tea.MouseButtonWheelDown:
		// Check if mouse is over preview pane
		if m.ShowPreview && m.isMouseOverPreview(msg) {
			// Mouse is over preview - scroll preview content
			m.PreviewScrollOffset += 3
			return m, nil
		}
		// Otherwise scroll card list
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, 1) // Move down one row
		} else {
			m.moveSelection(1)
			// In list view with preview: keep preview locked (don't update PreviewedIndex)
			// Preview only updates on keyboard nav, click, or spacebar pin
		}
		return m, nil

	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionRelease {
			return m.handleLeftClick(msg)
		}
	}

	return m, nil
}

// handleLeftClick processes left mouse button clicks (single and double-click)
func (m Model) handleLeftClick(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if len(m.FilteredCards) == 0 {
		return m, nil
	}

	// Calculate which card was clicked based on view mode
	clickedIndex := m.calculateClickedCardIndex(msg)

	// Invalid click (outside card area)
	if clickedIndex < 0 || clickedIndex >= len(m.FilteredCards) {
		return m, nil
	}

	// Double-click detection: same card clicked within 500ms
	const doubleClickThreshold = 500 * time.Millisecond
	now := time.Now()
	isDoubleClick := clickedIndex == m.LastClickIndex &&
		now.Sub(m.LastClickTime) < doubleClickThreshold

	if isDoubleClick {
		// Double-click: copy card to clipboard
		card := &m.FilteredCards[clickedIndex]
		m.LastClickIndex = -1
		m.LastClickTime = time.Time{}
		return m, copyToClipboard(card.Content)
	} else {
		// Single-click: select card and update preview
		m.SelectedIndex = clickedIndex
		m.PreviewedIndex = clickedIndex // Update preview on click
		m.PreviewScrollOffset = 0       // Reset preview scroll when changing card
		m.LastClickIndex = clickedIndex
		m.LastClickTime = now

		// Update scroll offset to keep selected card visible
		if m.ViewMode == ViewGrid {
			m.ensureGridSelectionVisible()
		} else {
			m.ensureListSelectionVisible()
		}
	}

	return m, nil
}

// calculateClickedCardIndex determines which card was clicked based on mouse position
func (m Model) calculateClickedCardIndex(msg tea.MouseMsg) int {
	// Header is 2 lines (title + border)
	// Status bar is 2 lines (border + hints)
	// Total overhead: 4 lines
	headerOffset := 2

	// Check if click is in valid area (below header, above status bar)
	// headerOffset is 2 (title + border), status bar is 2 (border + content)
	if msg.Y < headerOffset || msg.Y >= m.Height-2 {
		return -1
	}

	if m.ViewMode == ViewGrid {
		return m.calculateGridClickIndex(msg, headerOffset)
	} else if m.ViewMode == ViewTable {
		// Table view has 2 extra lines (header + separator) that we need to account for
		return m.calculateTableClickIndex(msg, headerOffset)
	} else {
		return m.calculateListClickIndex(msg, headerOffset)
	}
}

// calculateListClickIndex calculates which card was clicked in list view
func (m Model) calculateListClickIndex(msg tea.MouseMsg, headerOffset int) int {
	// Calculate available height for list view
	availableHeight := m.Height - 6 // Header + status bar

	if m.ShowPreview {
		// Split view: list gets half the height
		if availableHeight > 50 {
			availableHeight = availableHeight * 2 / 5 // 40% for list
		} else {
			availableHeight = availableHeight / 2 // 50% for list
		}
	}

	// Calculate which line was clicked (relative to list start)
	clickedLine := msg.Y - headerOffset

	// Check if click is within list area
	if clickedLine < 0 || clickedLine >= availableHeight {
		return -1
	}

	// Calculate the card index accounting for scroll offset
	clickedIndex := m.ScrollOffset + clickedLine

	// Validate bounds
	if clickedIndex >= len(m.FilteredCards) {
		return -1
	}

	return clickedIndex
}

// calculateTableClickIndex calculates which card was clicked in table view
func (m Model) calculateTableClickIndex(msg tea.MouseMsg, headerOffset int) int {
	// Table view has 2 extra lines at the top: header row + separator
	// These come before the data rows
	tableHeaderLines := 2

	// Calculate which line was clicked (relative to table start)
	clickedLine := msg.Y - headerOffset

	// Check if click is on the table header or separator (not a data row)
	if clickedLine < tableHeaderLines {
		return -1
	}

	// Subtract the table header lines to get the data row index
	dataRowIndex := clickedLine - tableHeaderLines

	// Calculate available height for table rows
	availableHeight := m.Height - 6 - tableHeaderLines // Header + status bar + table header

	// Check if click is within table data area
	if dataRowIndex < 0 || dataRowIndex >= availableHeight {
		return -1
	}

	// Calculate the card index accounting for scroll offset
	clickedIndex := m.ScrollOffset + dataRowIndex

	// Validate bounds
	if clickedIndex >= len(m.FilteredCards) {
		return -1
	}

	return clickedIndex
}

// calculateGridClickIndex calculates which card was clicked in grid view
func (m Model) calculateGridClickIndex(msg tea.MouseMsg, headerOffset int) int {
	// Calculate available dimensions for grid
	availableWidth := m.Width
	availableHeight := m.Height - 6
	gridWidth := availableWidth

	if m.ShowPreview {
		if m.Width > 120 {
			// Side-by-side: calculate actual grid width (matches renderGridWithPreview)
			separator := 2
			previewBorderPadding := 4
			actualAvailableWidth := m.Width - separator - previewBorderPadding
			gridWidth = actualAvailableWidth / 2
			if gridWidth < GridCardTotalWidth {
				gridWidth = GridCardTotalWidth
			}
			availableWidth = gridWidth

			// Check if click is within grid area (not in preview)
			if msg.X >= gridWidth {
				return -1 // Click is in preview area, not grid
			}
		} else {
			// Top/bottom: adjust height
			if availableHeight > 50 {
				availableHeight = availableHeight * 2 / 5 // Grid gets 40%
			} else {
				availableHeight = availableHeight / 2
			}
		}
	}

	// Calculate grid dimensions (MUST match rendering logic)
	cols := max(1, min(availableWidth/GridCardTotalWidth, GridMaxColumns))

	// Calculate which row and column were clicked
	clickedRow := (msg.Y - headerOffset) / GridCardTotalHeight
	clickedCol := msg.X / GridCardTotalWidth

	// Validate column
	if clickedCol < 0 || clickedCol >= cols {
		return -1
	}

	// Calculate the starting row based on scroll offset
	startRow := m.ScrollOffset / cols

	// Calculate actual card index
	actualRow := startRow + clickedRow
	clickedIndex := actualRow*cols + clickedCol

	// Validate bounds
	if clickedIndex < 0 || clickedIndex >= len(m.FilteredCards) {
		return -1
	}

	return clickedIndex
}

// ensureListSelectionVisible updates scroll offset to keep selected card visible in list view
func (m *Model) ensureListSelectionVisible() {
	visibleCount := m.getVisibleCardCount()

	if m.SelectedIndex < m.ScrollOffset {
		m.ScrollOffset = m.SelectedIndex
	}
	if m.SelectedIndex >= m.ScrollOffset+visibleCount {
		m.ScrollOffset = m.SelectedIndex - visibleCount + 1
	}
}

// ensureGridSelectionVisible updates scroll offset to keep selected card visible in grid view
func (m *Model) ensureGridSelectionVisible() {
	// Calculate grid dimensions (MUST match rendering logic)
	availableWidth := m.Width
	availableHeight := m.Height - 6

	if m.ShowPreview && m.Width > 120 {
		// Side-by-side: use same calculation as rendering
		maxGridWidth := (GridCardTotalWidth * GridMaxColumns) + 2
		availableWidth = min(m.Width/2, maxGridWidth)
	} else if m.ShowPreview {
		if availableHeight > 50 {
			availableHeight = availableHeight * 2 / 5
		} else {
			availableHeight = availableHeight / 2
		}
	}

	cols := max(1, min(availableWidth/GridCardTotalWidth, GridMaxColumns))
	visibleRows := max(1, availableHeight/GridCardTotalHeight)

	// Get current row and scroll row
	currentRow := m.SelectedIndex / cols
	scrollRow := m.ScrollOffset / cols

	// Adjust scroll if selection is out of view
	if currentRow < scrollRow {
		m.ScrollOffset = currentRow * cols
	}
	if currentRow >= scrollRow+visibleRows {
		m.ScrollOffset = (currentRow - visibleRows + 1) * cols
	}
}

// isMouseOverPreview checks if the mouse is positioned over the preview pane
func (m Model) isMouseOverPreview(msg tea.MouseMsg) bool {
	if !m.ShowPreview {
		return false
	}

	headerOffset := 2
	availableHeight := m.Height - 6 // Header + status bar

	if m.ViewMode == ViewGrid && m.Width > 120 {
		// Grid view side-by-side: preview is on right side
		maxGridWidth := (GridCardTotalWidth * GridMaxColumns) + 2
		actualCols := max(1, min((min(m.Width/2, maxGridWidth))/GridCardTotalWidth, GridMaxColumns))
		actualGridWidth := actualCols * GridCardTotalWidth
		// Check if mouse X is to the right of grid (accounting for margin)
		return msg.X > actualGridWidth+2
	} else {
		// Top/bottom layout (both grid and list view): preview is on bottom
		var contentHeight int
		if m.ViewMode == ViewGrid {
			if availableHeight > 50 {
				contentHeight = availableHeight * 2 / 5 // Grid gets 40%
			} else {
				contentHeight = availableHeight / 2
			}
		} else {
			// List view
			if availableHeight > 50 {
				contentHeight = availableHeight * 2 / 5 // List gets 40%
			} else {
				contentHeight = availableHeight / 2
			}
		}
		// Check if mouse Y is below the content area (in preview area)
		return msg.Y >= headerOffset+contentHeight
	}
}
