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
	// Handle mouse wheel scrolling
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		// In grid view with preview: check if mouse is over preview pane
		if m.ViewMode == ViewGrid && m.ShowPreview && m.Width > 120 {
			// Side-by-side layout: preview is on right side
			maxGridWidth := (GridCardTotalWidth * GridMaxColumns) + 2
			gridWidth := min(m.Width/2, maxGridWidth)
			if msg.X > gridWidth {
				// Mouse is over preview - scroll preview content
				m.PreviewScrollOffset = max(0, m.PreviewScrollOffset-3)
				return m, nil
			}
		}
		// Otherwise scroll card list
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, -1) // Move up one row
		} else {
			m.moveSelection(-1)
		}
		return m, nil

	case tea.MouseButtonWheelDown:
		// In grid view with preview: check if mouse is over preview pane
		if m.ViewMode == ViewGrid && m.ShowPreview && m.Width > 120 {
			// Side-by-side layout: preview is on right side
			maxGridWidth := (GridCardTotalWidth * GridMaxColumns) + 2
			gridWidth := min(m.Width/2, maxGridWidth)
			if msg.X > gridWidth {
				// Mouse is over preview - scroll preview content
				m.PreviewScrollOffset += 3
				return m, nil
			}
		}
		// Otherwise scroll card list
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, 1) // Move down one row
		} else {
			m.moveSelection(1)
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
	// Header is 3 lines (title + border)
	// Status bar is 2 lines (border + hints)
	// Total overhead: 5 lines
	headerOffset := 3

	// Check if click is in valid area (below header, above status bar)
	if msg.Y < headerOffset || msg.Y >= m.Height-2 {
		return -1
	}

	if m.ViewMode == ViewGrid {
		return m.calculateGridClickIndex(msg, headerOffset)
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

// calculateGridClickIndex calculates which card was clicked in grid view
func (m Model) calculateGridClickIndex(msg tea.MouseMsg, headerOffset int) int {
	// Calculate available dimensions for grid
	availableWidth := m.Width
	availableHeight := m.Height - 6

	if m.ShowPreview {
		if m.Width > 120 {
			// Side-by-side: use same calculation as rendering
			maxGridWidth := (GridCardTotalWidth * GridMaxColumns) + 2
			availableWidth = min(m.Width/2, maxGridWidth)
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
