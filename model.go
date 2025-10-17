package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// model.go - Model Initialization and Helpers
// Purpose: Create initial state and helper methods

// initialModel creates the initial application state
func initialModel() Model {
	return Model{
		Data:                nil, // Will be loaded asynchronously
		FilteredCards:       []Card{},
		CategoryMap:         make(map[string]Category),
		SelectedIndex:       0,
		PreviewedIndex:      0,
		PreviewScrollOffset: 0,
		ScrollOffset:        0,
		SearchQuery:         "",
		SelectedCategories:  make(map[string]bool),
		ViewMode:            ViewList,
		ShowPreview:         false, // Start with preview off for cleaner initial layout
		ShowHelp:            false,
		TemplateVars:        make(map[string]string),
		Width:               80,
		Height:              24,
		Error:               nil,
		LastClickIndex:      -1,
		LastClickTime:       time.Time{},
	}
}

// Init is called when the program starts (Bubbletea lifecycle)
func (m Model) Init() tea.Cmd {
	// Load data asynchronously
	return loadDataAsync(DefaultDataPath)
}

// buildCategoryMap creates a fast lookup map for categories
func (m *Model) buildCategoryMap() {
	if m.Data == nil {
		return
	}

	m.CategoryMap = make(map[string]Category)
	for _, cat := range m.Data.Categories {
		m.CategoryMap[cat.ID] = cat
	}
}

// updateFilteredCards applies search and category filters to the card list
func (m *Model) updateFilteredCards() {
	if m.Data == nil {
		m.FilteredCards = []Card{}
		return
	}

	cards := m.Data.Cards

	// Apply search filter
	if m.SearchQuery != "" {
		cards = searchCards(cards, m.SearchQuery)
	}

	// Apply category filter
	if len(m.SelectedCategories) > 0 {
		cards = filterByCategories(cards, m.SelectedCategories)
	}

	m.FilteredCards = cards

	// Adjust selected index if out of bounds
	if m.SelectedIndex >= len(m.FilteredCards) {
		m.SelectedIndex = max(0, len(m.FilteredCards)-1)
	}
}

// getSelectedCard returns the currently selected card, or nil if none
func (m *Model) getSelectedCard() *Card {
	if len(m.FilteredCards) == 0 || m.SelectedIndex < 0 || m.SelectedIndex >= len(m.FilteredCards) {
		return nil
	}
	return &m.FilteredCards[m.SelectedIndex]
}

// getPreviewedCard returns the card shown in preview pane, or nil if none
func (m *Model) getPreviewedCard() *Card {
	if len(m.FilteredCards) == 0 || m.PreviewedIndex < 0 || m.PreviewedIndex >= len(m.FilteredCards) {
		return nil
	}
	return &m.FilteredCards[m.PreviewedIndex]
}

// getCategoryForCard returns the category for a given card
func (m *Model) getCategoryForCard(card *Card) *Category {
	if card == nil {
		return nil
	}
	if cat, ok := m.CategoryMap[card.CategoryID]; ok {
		return &cat
	}
	return nil
}

// getVisibleCardCount returns the number of cards that can fit in the list view
func (m *Model) getVisibleCardCount() int {
	// Header (3 lines) + Status bar (1 line) + padding (2 lines)
	usedHeight := 6
	availableHeight := m.Height - usedHeight

	if m.ShowPreview {
		// Split view: use half the height for list
		availableHeight = availableHeight / 2
	}

	return max(1, availableHeight)
}

// moveSelection moves the cursor up or down
func (m *Model) moveSelection(delta int) {
	if len(m.FilteredCards) == 0 {
		return
	}

	m.SelectedIndex += delta

	// Clamp to valid range
	if m.SelectedIndex < 0 {
		m.SelectedIndex = 0
	}
	if m.SelectedIndex >= len(m.FilteredCards) {
		m.SelectedIndex = len(m.FilteredCards) - 1
	}

	// Update scroll offset to keep selection visible
	visibleCount := m.getVisibleCardCount()
	if m.SelectedIndex < m.ScrollOffset {
		m.ScrollOffset = m.SelectedIndex
	}
	if m.SelectedIndex >= m.ScrollOffset+visibleCount {
		m.ScrollOffset = m.SelectedIndex - visibleCount + 1
	}
}

// moveSelectionGrid moves the cursor in grid view (2D navigation)
func (m *Model) moveSelectionGrid(dx, dy int) {
	if len(m.FilteredCards) == 0 {
		return
	}

	// Calculate grid dimensions (account for preview pane if shown)
	availableWidth := m.Width
	if m.ShowPreview && m.Width > 120 {
		// Side-by-side layout: grid gets 60% of width
		availableWidth = m.Width * 3 / 5
	}
	cols := max(1, min(availableWidth/GridCardTotalWidth, GridMaxColumns))

	// Get current row and column
	currentRow := m.SelectedIndex / cols
	currentCol := m.SelectedIndex % cols

	// Calculate new position
	newRow := currentRow + dy
	newCol := currentCol + dx

	// Clamp column
	if newCol < 0 {
		newCol = 0
	}
	if newCol >= cols {
		newCol = cols - 1
	}

	// Clamp row
	maxRow := (len(m.FilteredCards) - 1) / cols
	if newRow < 0 {
		newRow = 0
	}
	if newRow > maxRow {
		newRow = maxRow
	}

	// Calculate new index
	newIndex := newRow*cols + newCol

	// Clamp to valid card index
	if newIndex >= len(m.FilteredCards) {
		newIndex = len(m.FilteredCards) - 1
	}
	if newIndex < 0 {
		newIndex = 0
	}

	m.SelectedIndex = newIndex

	// Update scroll offset (in grid, we scroll by rows)
	availableHeight := m.Height - 6
	if m.ShowPreview && m.Width <= 120 {
		// Top/bottom layout: adjust available height
		if availableHeight > 50 {
			availableHeight = availableHeight * 2 / 5 // Grid gets 40%
		} else {
			availableHeight = availableHeight / 2
		}
	}
	visibleRows := max(1, availableHeight/GridCardTotalHeight)

	scrollRow := m.ScrollOffset / cols
	if newRow < scrollRow {
		m.ScrollOffset = newRow * cols
	}
	if newRow >= scrollRow+visibleRows {
		m.ScrollOffset = (newRow - visibleRows + 1) * cols
	}

	// Clamp scroll offset to valid range
	maxScrollOffset := max(0, len(m.FilteredCards)-1)
	if m.ScrollOffset > maxScrollOffset {
		m.ScrollOffset = maxScrollOffset
	}
	if m.ScrollOffset < 0 {
		m.ScrollOffset = 0
	}
}

// toggleCategory toggles a category filter on/off
func (m *Model) toggleCategory(categoryID string) {
	if m.SelectedCategories[categoryID] {
		delete(m.SelectedCategories, categoryID)
	} else {
		m.SelectedCategories[categoryID] = true
	}
	m.updateFilteredCards()
}

// clearFilters resets all filters
func (m *Model) clearFilters() {
	m.SearchQuery = ""
	m.SelectedCategories = make(map[string]bool)
	m.updateFilteredCards()
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// saveNewCard creates a new card and saves it to disk
func (m *Model) saveNewCard() tea.Cmd {
	// Validate input
	if m.NewCardTitle == "" || m.NewCardContent == "" {
		return nil // TODO: Show error message
	}

	// Create the new card
	return func() tea.Msg {
		newCard := Card{
			ID:         generateCardID(),
			Title:      m.NewCardTitle,
			Content:    m.NewCardContent,
			CategoryID: m.NewCardCategoryID,
			CreatedAt:  time.Now().UnixMilli(),
			UpdatedAt:  time.Now().UnixMilli(),
		}

		// Add to data
		m.Data.Cards = append(m.Data.Cards, newCard)

		// Save to disk
		if err := SaveData(DefaultDataPath, m.Data); err != nil {
			return cardSaveErrorMsg{err: err}
		}

		return cardSavedMsg{card: &newCard}
	}
}
