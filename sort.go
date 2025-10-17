package main

import (
	"sort"
	"strings"
	"time"
)

// sort.go - Card Sorting Functions
// Purpose: Sort cards by different columns for table view

// sortCards sorts a slice of cards based on the specified column and direction
func sortCards(cards []Card, categoryMap map[string]Category, column, direction string) []Card {
	// Make a copy to avoid modifying the original
	sorted := make([]Card, len(cards))
	copy(sorted, cards)

	// Define the comparison function based on column
	switch column {
	case "title":
		sort.SliceStable(sorted, func(i, j int) bool {
			result := strings.ToLower(sorted[i].Title) < strings.ToLower(sorted[j].Title)
			if direction == "desc" {
				return !result
			}
			return result
		})

	case "category":
		sort.SliceStable(sorted, func(i, j int) bool {
			// Get category names for comparison
			catI := categoryMap[sorted[i].CategoryID]
			catJ := categoryMap[sorted[j].CategoryID]
			result := strings.ToLower(catI.Name) < strings.ToLower(catJ.Name)
			if direction == "desc" {
				return !result
			}
			return result
		})

	case "created":
		sort.SliceStable(sorted, func(i, j int) bool {
			result := sorted[i].CreatedAt < sorted[j].CreatedAt
			if direction == "desc" {
				return !result
			}
			return result
		})

	case "updated":
		sort.SliceStable(sorted, func(i, j int) bool {
			result := sorted[i].UpdatedAt < sorted[j].UpdatedAt
			if direction == "desc" {
				return !result
			}
			return result
		})

	default:
		// Default to title sorting
		sort.SliceStable(sorted, func(i, j int) bool {
			result := strings.ToLower(sorted[i].Title) < strings.ToLower(sorted[j].Title)
			if direction == "desc" {
				return !result
			}
			return result
		})
	}

	return sorted
}

// formatDate formats Unix millisecond timestamp to human-readable date
func formatDate(unixMillis int64) string {
	if unixMillis == 0 {
		return "N/A"
	}
	t := time.Unix(unixMillis/1000, 0)
	return t.Format("2006-01-02")
}

// formatDateTime formats Unix millisecond timestamp to human-readable date and time
func formatDateTime(unixMillis int64) string {
	if unixMillis == 0 {
		return "N/A"
	}
	t := time.Unix(unixMillis/1000, 0)
	return t.Format("2006-01-02 15:04")
}

// getSortIndicator returns the sort direction indicator (↑ or ↓)
func getSortIndicator(column, currentColumn, direction string) string {
	if column != currentColumn {
		return ""
	}
	if direction == "asc" {
		return " ↑"
	}
	return " ↓"
}
