package main

import "time"

// types.go - Core Data Structures
// Purpose: Define all data models and application state

// Card represents a CellBlocks card (command, prompt, snippet, etc.)
type Card struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID string `json:"categoryId"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
	ImageID    string `json:"imageId,omitempty"`
}

// Category represents a card category with color theming
type Category struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Color            string `json:"color"`
	Hidden           bool   `json:"hidden,omitempty"`
	ParentCategoryID string `json:"parentCategoryId,omitempty"`
}

// CellBlocksData is the root structure matching cellblocks-data.json
type CellBlocksData struct {
	Version    string     `json:"version,omitempty"`
	ExportedAt string     `json:"exportedAt,omitempty"`
	Cards      []Card     `json:"cards"`
	Categories []Category `json:"categories"`
}

// ViewMode defines the current view layout
type ViewMode int

const (
	ViewList ViewMode = iota
	ViewGrid
	ViewDetail
)

// Model is the main application state (Bubbletea Model)
type Model struct {
	// Data
	Data          *CellBlocksData
	FilteredCards []Card
	CategoryMap   map[string]Category // Fast category lookup by ID

	// UI State
	SelectedIndex      int
	PreviewedIndex     int // Card shown in preview (only updates on click)
	PreviewScrollOffset int // Scroll position within preview content
	ScrollOffset       int
	SearchQuery        string
	SelectedCategories map[string]bool // Set of selected category IDs

	// View mode
	ViewMode    ViewMode
	ShowPreview bool
	ShowHelp    bool

	// Template editing
	TemplateVars map[string]string

	// Terminal size
	Width  int
	Height int

	// Error state
	Error error

	// Click tracking for double-click detection
	LastClickIndex int
	LastClickTime  time.Time
}

// Messages for Bubbletea update loop

// dataLoadedMsg is sent when card data is successfully loaded
type dataLoadedMsg struct {
	data *CellBlocksData
}

// dataLoadErrorMsg is sent when data loading fails
type dataLoadErrorMsg struct {
	err error
}

// cardCopiedMsg is sent when a card is copied to clipboard
type cardCopiedMsg struct {
	cardTitle string
}

// copyErrorMsg is sent when clipboard copy fails
type copyErrorMsg struct {
	err error
}
