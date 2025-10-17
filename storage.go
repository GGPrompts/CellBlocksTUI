package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// generateCardID creates a unique ID for a new card
func generateCardID() string {
	// Generate 16 random bytes (128 bits)
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID if random fails
		return fmt.Sprintf("card_%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// GetFileModTime returns the modification time of the data file
func GetFileModTime(path string) (time.Time, error) {
	fullPath := expandPath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// checkFileChanges checks if the data file has been modified and reloads if needed
func checkFileChanges(path string, lastModTime time.Time, currentCardCount int) tea.Cmd {
	return func() tea.Msg {
		// Check file modification time
		modTime, err := GetFileModTime(path)
		if err != nil {
			// File doesn't exist or can't be accessed - ignore
			return nil
		}

		// If file hasn't changed, return tick to check again later
		if !modTime.After(lastModTime) {
			return tickMsg{}
		}

		// File has changed - reload data
		data, err := LoadData(path)
		if err != nil {
			// Failed to load - return tick to try again later
			return tickMsg{}
		}

		// Calculate how many new cards were added
		newCards := len(data.Cards) - currentCardCount

		return fileChangedMsg{
			data:     data,
			newCards: newCards,
		}
	}
}

// startFileTicker starts a periodic ticker to check for file changes
func startFileTicker() tea.Cmd {
	return tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
