package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// main.go - Application Entry Point
// Purpose: ONLY contains the main() function
// Rule: Never add business logic to this file. Keep it minimal.

func main() {
	// Create program with options
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse clicks (without constant hover events)
	}

	p := tea.NewProgram(
		initialModel(),
		opts...,
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
