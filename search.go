package main

import (
	"strings"
)

// search.go - Search and Filtering Engine
// Purpose: Full-text search and category filtering

// searchCards filters cards by search query (title and content)
func searchCards(cards []Card, query string) []Card {
	if query == "" {
		return cards
	}

	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return cards
	}

	var results []Card
	for _, card := range cards {
		if matchesSearch(card, query) {
			results = append(results, card)
		}
	}

	return results
}

// matchesSearch checks if a card matches the search query
func matchesSearch(card Card, query string) bool {
	// Search in title
	if strings.Contains(strings.ToLower(card.Title), query) {
		return true
	}

	// Search in content
	if strings.Contains(strings.ToLower(card.Content), query) {
		return true
	}

	return false
}

// filterByCategories filters cards by selected categories
func filterByCategories(cards []Card, selectedCategories map[string]bool) []Card {
	if len(selectedCategories) == 0 {
		return cards
	}

	var results []Card
	for _, card := range cards {
		if selectedCategories[card.CategoryID] {
			results = append(results, card)
		}
	}

	return results
}
