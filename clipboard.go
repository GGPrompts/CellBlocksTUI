package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

// clipboard.go - Clipboard Integration
// Purpose: Platform-specific clipboard operations

// detectPlatform returns the clipboard command for the current platform
func detectPlatform() (string, []string) {
	// Check if running in Termux
	if isTermux() {
		return "termux-clipboard-set", []string{}
	}

	// Platform-specific commands
	switch runtime.GOOS {
	case "darwin":
		return "pbcopy", []string{}
	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			return "xclip", []string{"-selection", "clipboard"}
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			return "xsel", []string{"--clipboard", "--input"}
		}
	case "windows":
		return "clip.exe", []string{}
	}

	return "", nil
}

// isTermux checks if running in Termux environment
func isTermux() bool {
	// Termux sets the PREFIX environment variable
	if prefix := os.Getenv("PREFIX"); prefix != "" {
		return true
	}
	// Also check if termux-clipboard-set exists
	if _, err := exec.LookPath("termux-clipboard-set"); err == nil {
		return true
	}
	return false
}

// copyToClipboardSync performs the actual clipboard copy operation
func copyToClipboardSync(text string) error {
	cmd, args := detectPlatform()
	if cmd == "" {
		return fmt.Errorf("clipboard not supported on this platform")
	}

	// Create command
	command := exec.Command(cmd, args...)

	// Write text to stdin
	stdin, err := command.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	// Start command
	if err := command.Start(); err != nil {
		return fmt.Errorf("failed to start clipboard command: %w", err)
	}

	// Write content
	if _, err := stdin.Write([]byte(text)); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}
	stdin.Close()

	// Wait for command to finish
	if err := command.Wait(); err != nil {
		return fmt.Errorf("clipboard command failed: %w", err)
	}

	return nil
}

// copyToClipboard is the async Bubbletea command wrapper
func copyToClipboard(text string) tea.Cmd {
	return func() tea.Msg {
		err := copyToClipboardSync(text)
		if err != nil {
			return copyErrorMsg{err: err}
		}
		return cardCopiedMsg{cardTitle: "Card copied to clipboard"}
	}
}
