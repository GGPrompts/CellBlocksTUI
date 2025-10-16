package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// storage.go - Data Loading and Persistence
// Purpose: Read/write cellblocks-data.json

const (
	// DefaultDataPath is the shared data file location
	DefaultDataPath = "~/projects/CellBlocks/data/cellblocks-data.json"
)

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path // Return as-is if we can't get home dir
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// LoadData reads and parses the cellblocks-data.json file
func LoadData(path string) (*CellBlocksData, error) {
	// Expand ~ to home directory
	fullPath := expandPath(path)

	// Read file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	// Parse JSON
	var data CellBlocksData
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &data, nil
}

// loadDataAsync loads data in the background and sends a Bubbletea message
func loadDataAsync(path string) tea.Cmd {
	return func() tea.Msg {
		data, err := LoadData(path)
		if err != nil {
			return dataLoadErrorMsg{err: err}
		}
		return dataLoadedMsg{data: data}
	}
}

// SaveData writes the CellBlocksData to disk (for future editing features)
func SaveData(path string, data *CellBlocksData) error {
	fullPath := expandPath(path)

	// Marshal to JSON with indentation
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
}

// FileExists checks if the data file exists
func FileExists(path string) bool {
	fullPath := expandPath(path)
	_, err := os.Stat(fullPath)
	return err == nil
}
