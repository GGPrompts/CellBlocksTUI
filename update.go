package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// update.go - Main Update Loop
// Purpose: Handle all Bubbletea messages and update model state

// Update is the main update function (Bubbletea lifecycle)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Window resize
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	// Data loaded successfully
	case dataLoadedMsg:
		m.Data = msg.data
		m.buildCategoryMap()
		m.updateFilteredCards()
		return m, nil

	// Data loading failed
	case dataLoadErrorMsg:
		m.Error = msg.err
		return m, nil

	// Keyboard events
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	// Mouse events
	case tea.MouseMsg:
		return m.handleMouseEventEnhanced(msg)
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global shortcuts
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "?":
		m.ShowHelp = !m.ShowHelp
		return m, nil

	case "p":
		// Toggle preview pane
		m.ShowPreview = !m.ShowPreview
		// When enabling preview, sync previewed card to current selection
		if m.ShowPreview {
			m.PreviewedIndex = m.SelectedIndex
			m.PreviewScrollOffset = 0
		}
		return m, nil

	case "g":
		// Toggle between list and grid view
		if m.ViewMode == ViewGrid {
			m.ViewMode = ViewList
		} else {
			m.ViewMode = ViewGrid
			// Reset scroll offset when switching to grid view
			m.ScrollOffset = 0
		}
		return m, nil

	case "esc":
		if m.ShowHelp {
			m.ShowHelp = false
			return m, nil
		}
		if m.SearchQuery != "" {
			m.SearchQuery = ""
			m.updateFilteredCards()
			return m, nil
		}
		return m, nil
	}

	// If help is shown, don't process other keys
	if m.ShowHelp {
		return m, nil
	}

	// Navigation
	switch msg.String() {
	case "up", "k":
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, -1) // Move up one row
		} else {
			m.moveSelection(-1)
			// In list view, update preview as you navigate
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "down", "j":
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, 1) // Move down one row
		} else {
			m.moveSelection(1)
			// In list view, update preview as you navigate
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "left", "h":
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(-1, 0) // Move left
		}
		return m, nil

	case "right", "l":
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(1, 0) // Move right
		}
		return m, nil

	case "shift+up":
		// Scroll preview up (in grid view with preview)
		if m.ViewMode == ViewGrid && m.ShowPreview {
			m.PreviewScrollOffset = max(0, m.PreviewScrollOffset-3)
		}
		return m, nil

	case "shift+down":
		// Scroll preview down (in grid view with preview)
		if m.ViewMode == ViewGrid && m.ShowPreview {
			m.PreviewScrollOffset += 3
		}
		return m, nil

	case "pageup":
		m.moveSelection(-m.getVisibleCardCount())
		// In list view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "pagedown":
		m.moveSelection(m.getVisibleCardCount())
		// In list view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "home":
		m.SelectedIndex = 0
		m.ScrollOffset = 0
		// In list view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "end":
		m.SelectedIndex = max(0, len(m.FilteredCards)-1)
		visibleCount := m.getVisibleCardCount()
		m.ScrollOffset = max(0, len(m.FilteredCards)-visibleCount)
		// In list view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "enter", "c":
		// Copy selected card to clipboard
		card := m.getSelectedCard()
		if card != nil {
			return m, copyToClipboard(card.Content)
		}
		return m, nil

	case " ": // Spacebar
		// In grid view: pin current selection to preview
		// In list view: toggle preview (existing behavior handled above by "p")
		if m.ViewMode == ViewGrid && m.ShowPreview {
			m.PreviewedIndex = m.SelectedIndex
			m.PreviewScrollOffset = 0 // Reset scroll when pinning new card
		}
		return m, nil

	case "/":
		// Toggle search mode (for now, just clear filters)
		m.clearFilters()
		return m, nil
	}

	// Search input (simple character by character for now)
	if len(msg.String()) == 1 {
		m.SearchQuery += msg.String()
		m.updateFilteredCards()
		return m, nil
	}

	if msg.Type == tea.KeyBackspace && len(m.SearchQuery) > 0 {
		m.SearchQuery = m.SearchQuery[:len(m.SearchQuery)-1]
		m.updateFilteredCards()
		return m, nil
	}

	return m, nil
}

// handleMouseEvent is replaced by handleMouseEventEnhanced in update_mouse.go
// Kept here as reference for the old basic implementation
// func (m Model) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
// 	switch msg.Type {
// 	case tea.MouseWheelUp:
// 		m.moveSelection(-1)
// 		return m, nil
// 	case tea.MouseWheelDown:
// 		m.moveSelection(1)
// 		return m, nil
// 	}
// 	return m, nil
// }
