package main

import (
	"fmt"
	"time"

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
		// Get initial file modification time
		if modTime, err := GetFileModTime(DefaultDataPath); err == nil {
			m.LastFileModTime = modTime
		}
		// Start the file change checker
		return m, startFileTicker()

	// Data loading failed
	case dataLoadErrorMsg:
		m.Error = msg.err
		return m, nil

	// Card saved successfully
	case cardSavedMsg:
		// Refresh the filtered cards to include the new card
		m.updateFilteredCards()
		// Update file modification time after save
		if modTime, err := GetFileModTime(DefaultDataPath); err == nil {
			m.LastFileModTime = modTime
		}
		// Return to list view
		m.ViewMode = ViewList
		// Jump to the new card (it should be at the end)
		if len(m.FilteredCards) > 0 {
			m.SelectedIndex = len(m.FilteredCards) - 1
			m.PreviewedIndex = m.SelectedIndex
			// Scroll to bottom
			visibleCount := m.getVisibleCardCount()
			m.ScrollOffset = max(0, len(m.FilteredCards)-visibleCount)
		}
		return m, nil

	// Card save failed
	case cardSaveErrorMsg:
		m.Error = msg.err
		return m, nil

	// Periodic tick to check for file changes
	case tickMsg:
		if m.Data != nil {
			return m, checkFileChanges(DefaultDataPath, m.LastFileModTime, len(m.Data.Cards))
		}
		return m, startFileTicker()

	// File changed externally - reload data
	case fileChangedMsg:
		m.Data = msg.data
		m.buildCategoryMap()
		m.updateFilteredCards()
		// Update file modification time
		if modTime, err := GetFileModTime(DefaultDataPath); err == nil {
			m.LastFileModTime = modTime
		}
		// Show notification if new cards were added
		if msg.newCards > 0 {
			m.ReloadMessage = fmt.Sprintf("✨ %d new card(s) detected!", msg.newCards)
			m.ReloadMessageTime = time.Now()
		} else if msg.newCards < 0 {
			m.ReloadMessage = fmt.Sprintf("🔄 Data reloaded (%d card(s) removed)", -msg.newCards)
			m.ReloadMessageTime = time.Now()
		} else {
			m.ReloadMessage = "🔄 Data reloaded"
			m.ReloadMessageTime = time.Now()
		}
		// Continue checking for changes
		return m, startFileTicker()

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
		// Reset scroll offset in grid view when toggling preview
		// because available height changes and cards might be off-screen
		if m.ViewMode == ViewGrid {
			m.ScrollOffset = 0
			// Also reset selection to ensure it's visible
			if m.SelectedIndex >= len(m.FilteredCards) {
				m.SelectedIndex = max(0, len(m.FilteredCards)-1)
			}
		}
		return m, nil

	case "m":
		// Toggle markdown rendering (works in list/grid view with preview, and detail view)
		if m.ViewMode == ViewList || m.ViewMode == ViewGrid || m.ViewMode == ViewDetail {
			m.UseMarkdownRender = !m.UseMarkdownRender
		}
		return m, nil

	case "g":
		// Cycle through list → grid → table → list
		switch m.ViewMode {
		case ViewList:
			m.ViewMode = ViewGrid
			m.ScrollOffset = 0
		case ViewGrid:
			m.ViewMode = ViewTable
			m.ScrollOffset = 0
		case ViewTable:
			m.ViewMode = ViewList
			m.ScrollOffset = 0
		default:
			m.ViewMode = ViewList
		}
		return m, nil

	case "f":
		// Open category filter screen
		if m.ViewMode == ViewList || m.ViewMode == ViewGrid || m.ViewMode == ViewTable {
			m.ViewMode = ViewCategoryFilter
			m.FilterCursorIndex = 0
		}
		return m, nil

	case "n":
		// Open card creation screen
		if m.ViewMode == ViewList || m.ViewMode == ViewGrid || m.ViewMode == ViewTable {
			m.ViewMode = ViewCardCreate
			m.CreateFormField = 0
			m.NewCardTitle = ""
			m.NewCardContent = ""
			// Default to first category if available
			if m.Data != nil && len(m.Data.Categories) > 0 {
				m.NewCardCategoryID = m.Data.Categories[0].ID
			}
		}
		return m, nil

	case "1":
		// Sort by title (table view only)
		if m.ViewMode == ViewTable {
			if m.SortColumn == "title" {
				// Toggle direction
				if m.SortDirection == "asc" {
					m.SortDirection = "desc"
				} else {
					m.SortDirection = "asc"
				}
			} else {
				m.SortColumn = "title"
				m.SortDirection = "asc"
			}
		}
		return m, nil

	case "2":
		// Sort by category (table view only)
		if m.ViewMode == ViewTable {
			if m.SortColumn == "category" {
				// Toggle direction
				if m.SortDirection == "asc" {
					m.SortDirection = "desc"
				} else {
					m.SortDirection = "asc"
				}
			} else {
				m.SortColumn = "category"
				m.SortDirection = "asc"
			}
		}
		return m, nil

	case "3":
		// Sort by created date (table view only)
		if m.ViewMode == ViewTable {
			if m.SortColumn == "created" {
				// Toggle direction
				if m.SortDirection == "asc" {
					m.SortDirection = "desc"
				} else {
					m.SortDirection = "asc"
				}
			} else {
				m.SortColumn = "created"
				m.SortDirection = "asc"
			}
		}
		return m, nil

	case "4":
		// Sort by updated date (table view only)
		if m.ViewMode == ViewTable {
			if m.SortColumn == "updated" {
				// Toggle direction
				if m.SortDirection == "asc" {
					m.SortDirection = "desc"
				} else {
					m.SortDirection = "asc"
				}
			} else {
				m.SortColumn = "updated"
				m.SortDirection = "asc"
			}
		}
		return m, nil

	case "esc":
		if m.ShowHelp {
			m.ShowHelp = false
			return m, nil
		}
		// Exit special screens back to main view
		if m.ViewMode == ViewCategoryFilter || m.ViewMode == ViewCardCreate || m.ViewMode == ViewDetail {
			// Reset detail view state when exiting detail mode
			m.DetailScrollOffset = 0
			m.ShowTemplateForm = false
			// Return to list view
			m.ViewMode = ViewList
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

	// Category filter screen handlers
	if m.ViewMode == ViewCategoryFilter {
		return m.handleCategoryFilterInput(msg)
	}

	// Card creation screen handlers
	if m.ViewMode == ViewCardCreate {
		return m.handleCardCreateInput(msg)
	}

	// Detail view handlers
	if m.ViewMode == ViewDetail {
		return m.handleDetailViewInput(msg)
	}

	// Navigation
	switch msg.String() {
	case "up", "k":
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, -1) // Move up one row
		} else {
			m.moveSelection(-1)
			// In list/table view, update preview as you navigate
			if m.ViewMode == ViewList {
				m.PreviewedIndex = m.SelectedIndex
			}
		}
		return m, nil

	case "down", "j":
		if m.ViewMode == ViewGrid {
			m.moveSelectionGrid(0, 1) // Move down one row
		} else {
			m.moveSelection(1)
			// In list/table view, update preview as you navigate
			if m.ViewMode == ViewList {
				m.PreviewedIndex = m.SelectedIndex
			}
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
		// Scroll preview up (when preview is shown)
		if m.ShowPreview {
			m.PreviewScrollOffset = max(0, m.PreviewScrollOffset-3)
		}
		return m, nil

	case "shift+down":
		// Scroll preview down (when preview is shown)
		if m.ShowPreview {
			m.PreviewScrollOffset += 3
		}
		return m, nil

	case "pageup":
		m.moveSelection(-m.getVisibleCardCount())
		// In list/table view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "pagedown":
		m.moveSelection(m.getVisibleCardCount())
		// In list/table view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "home":
		m.SelectedIndex = 0
		m.ScrollOffset = 0
		// In list/table view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "end":
		m.SelectedIndex = max(0, len(m.FilteredCards)-1)
		visibleCount := m.getVisibleCardCount()
		m.ScrollOffset = max(0, len(m.FilteredCards)-visibleCount)
		// In list/table view, update preview
		if m.ViewMode == ViewList {
			m.PreviewedIndex = m.SelectedIndex
		}
		return m, nil

	case "enter":
		// Enter detail view for selected card
		card := m.getSelectedCard()
		if card != nil {
			m.ViewMode = ViewDetail
			m.DetailScrollOffset = 0
			// Detect template variables
			m.DetectedVars = ExtractVariables(card.Content)
			// Initialize template vars if we have detected vars
			if len(m.DetectedVars) > 0 && m.TemplateVars == nil {
				m.TemplateVars = make(map[string]string)
			}
			// Auto-show template form if variables detected
			m.ShowTemplateForm = len(m.DetectedVars) > 0
		}
		return m, nil

	case "c":
		// Copy selected card to clipboard
		card := m.getSelectedCard()
		if card != nil {
			return m, copyToClipboard(card.Content)
		}
		return m, nil

	case "d":
		// Also enter detail view (alternative to Enter)
		card := m.getSelectedCard()
		if card != nil {
			m.ViewMode = ViewDetail
			m.DetailScrollOffset = 0
			// Detect template variables
			m.DetectedVars = ExtractVariables(card.Content)
			// Initialize template vars if we have detected vars
			if len(m.DetectedVars) > 0 && m.TemplateVars == nil {
				m.TemplateVars = make(map[string]string)
			}
			// Auto-show template form if variables detected
			m.ShowTemplateForm = len(m.DetectedVars) > 0
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

// handleCategoryFilterInput processes input in category filter screen
func (m Model) handleCategoryFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.Data == nil {
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		if m.FilterCursorIndex > 0 {
			m.FilterCursorIndex--
		}
		return m, nil

	case "down", "j":
		if m.FilterCursorIndex < len(m.Data.Categories)-1 {
			m.FilterCursorIndex++
		}
		return m, nil

	case "enter", " ":
		// Toggle selected category
		if m.FilterCursorIndex >= 0 && m.FilterCursorIndex < len(m.Data.Categories) {
			categoryID := m.Data.Categories[m.FilterCursorIndex].ID
			m.toggleCategory(categoryID)
		}
		return m, nil

	case "a":
		// Select all categories
		for _, cat := range m.Data.Categories {
			m.SelectedCategories[cat.ID] = true
		}
		m.updateFilteredCards()
		return m, nil

	case "c":
		// Clear all filters
		m.clearFilters()
		return m, nil
	}

	return m, nil
}

// handleCardCreateInput processes input in card creation screen
func (m Model) handleCardCreateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		// Move to next field
		m.CreateFormField = (m.CreateFormField + 1) % 3
		return m, nil

	case "shift+tab":
		// Move to previous field
		m.CreateFormField = (m.CreateFormField + 2) % 3 // +2 is same as -1 mod 3
		return m, nil

	case "ctrl+s", "ctrl+enter":
		// Save card
		return m, m.saveNewCard()

	case "backspace":
		// Delete character from current field
		switch m.CreateFormField {
		case 0: // Title
			if len(m.NewCardTitle) > 0 {
				m.NewCardTitle = m.NewCardTitle[:len(m.NewCardTitle)-1]
			}
		case 1: // Content
			if len(m.NewCardContent) > 0 {
				m.NewCardContent = m.NewCardContent[:len(m.NewCardContent)-1]
			}
		}
		return m, nil

	case "enter":
		// In content field, add newline
		if m.CreateFormField == 1 {
			m.NewCardContent += "\n"
		}
		return m, nil

	case "up", "k":
		// In category field, move to previous category
		if m.CreateFormField == 2 && m.Data != nil {
			for i, cat := range m.Data.Categories {
				if cat.ID == m.NewCardCategoryID && i > 0 {
					m.NewCardCategoryID = m.Data.Categories[i-1].ID
					break
				}
			}
		}
		return m, nil

	case "down", "j":
		// In category field, move to next category
		if m.CreateFormField == 2 && m.Data != nil {
			for i, cat := range m.Data.Categories {
				if cat.ID == m.NewCardCategoryID && i < len(m.Data.Categories)-1 {
					m.NewCardCategoryID = m.Data.Categories[i+1].ID
					break
				}
			}
		}
		return m, nil
	}

	// Type characters into current field
	if len(msg.String()) == 1 {
		switch m.CreateFormField {
		case 0: // Title
			m.NewCardTitle += msg.String()
		case 1: // Content
			m.NewCardContent += msg.String()
		}
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

// handleDetailViewInput processes input in detail view mode
func (m Model) handleDetailViewInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	card := m.getSelectedCard()
	if card == nil {
		return m, nil
	}

	switch msg.String() {
	case "up", "k":
		// Scroll content up
		m.DetailScrollOffset = max(0, m.DetailScrollOffset-1)
		return m, nil

	case "down", "j":
		// Scroll content down
		m.DetailScrollOffset++
		return m, nil

	case "pageup":
		// Scroll up by page
		m.DetailScrollOffset = max(0, m.DetailScrollOffset-10)
		return m, nil

	case "pagedown":
		// Scroll down by page
		m.DetailScrollOffset += 10
		return m, nil

	case "t":
		// Toggle template form visibility
		if len(m.DetectedVars) > 0 {
			m.ShowTemplateForm = !m.ShowTemplateForm
			if m.ShowTemplateForm {
				m.TemplateFormField = 0 // Reset to first field
			}
		}
		return m, nil

	case "c":
		// Copy card content (or filled template if form is shown)
		var contentToCopy string
		if m.ShowTemplateForm && len(m.DetectedVars) > 0 {
			// Copy filled template
			contentToCopy = FillTemplate(card.Content, m.TemplateVars)
		} else {
			// Copy raw content
			contentToCopy = card.Content
		}
		return m, copyToClipboard(contentToCopy)

	case "enter":
		// Copy filled template (if template form is shown)
		if m.ShowTemplateForm && len(m.DetectedVars) > 0 {
			contentToCopy := FillTemplate(card.Content, m.TemplateVars)
			return m, copyToClipboard(contentToCopy)
		}
		// Otherwise, just copy raw content
		return m, copyToClipboard(card.Content)

	case "tab":
		// Navigate to next template field (if template form is shown)
		if m.ShowTemplateForm && len(m.DetectedVars) > 0 {
			m.TemplateFormField = (m.TemplateFormField + 1) % len(m.DetectedVars)
		}
		return m, nil

	case "shift+tab":
		// Navigate to previous template field
		if m.ShowTemplateForm && len(m.DetectedVars) > 0 {
			m.TemplateFormField = (m.TemplateFormField - 1 + len(m.DetectedVars)) % len(m.DetectedVars)
		}
		return m, nil

	case "backspace":
		// Delete character from current template field
		if m.ShowTemplateForm && len(m.DetectedVars) > 0 && m.TemplateFormField < len(m.DetectedVars) {
			varName := m.DetectedVars[m.TemplateFormField]
			if len(m.TemplateVars[varName]) > 0 {
				m.TemplateVars[varName] = m.TemplateVars[varName][:len(m.TemplateVars[varName])-1]
			}
		}
		return m, nil
	}

	// Type characters into current template field
	if m.ShowTemplateForm && len(m.DetectedVars) > 0 && len(msg.String()) == 1 {
		if m.TemplateFormField < len(m.DetectedVars) {
			varName := m.DetectedVars[m.TemplateFormField]
			m.TemplateVars[varName] += msg.String()
		}
		return m, nil
	}

	return m, nil
}
