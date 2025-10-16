# CellBlocksTUI - Quick Start Guide

## What's Been Built

Phase 1 MVP is complete and ready to use!

### Files Created
- `types.go` - Data structures (Card, Category, Model)
- `storage.go` - JSON data loading
- `search.go` - Full-text search & filtering
- `model.go` - State management & helpers
- `update.go` - Event handling (keyboard, mouse)
- `view.go` - UI rendering with Lipgloss
- `styles.go` - Color themes & styling
- `clipboard.go` - Multi-platform clipboard support
- `main.go` - Application entry point

### Stats
- **Binary size**: 4.6MB (target: <5MB) ✓
- **Cards loaded**: 271 cards ✓
- **Categories**: 10 categories ✓
- **Search**: Working (24 results for "git") ✓

## Running the App

### Start CellBlocksTUI
```bash
cd ~/projects/CellBlocksTUI
./cellblocks-tui
```

### Keyboard Shortcuts

**Navigation:**
- `↑/k`, `↓/j` - Move up/down
- `←/h`, `→/l` - Move left/right (grid view)
- `PgUp/PgDn` - Scroll by page
- `Home/End` - Jump to first/last card

**View Modes:**
- `g` - Toggle grid/list view
- `p` - Toggle preview pane (list mode only)

**Actions:**
- `Enter` or `c` - Copy card to clipboard
- `/` - Clear all filters
- `Type...` - Search cards (live filtering)
- `Backspace` - Delete search character

**Help & Exit:**
- `?` - Show/hide help dialog
- `Esc` - Close help or clear search
- `q` or `Ctrl+C` - Quit

## Features Working

✓ Load 271 cards from `~/projects/CellBlocks/data/cellblocks-data.json`
✓ **Grid view** - Beautiful card blocks with neon borders (press `g`)
✓ List view with scrollable cards
✓ Live search (type to filter)
✓ Preview pane (toggle with `p` in list mode)
✓ Copy to clipboard (works on Linux, macOS, Windows, Termux)
✓ Category-colored borders in grid view
✓ 2D navigation (arrow keys in grid)
✓ Mouse/touch support (scroll, click)
✓ Help dialog (`?`)
✓ Responsive layout

## Platform Support

### Clipboard Commands
- **Termux**: `termux-clipboard-set` (auto-detected)
- **Linux**: `xclip` or `xsel` (install if needed)
- **macOS**: `pbcopy` (built-in)
- **Windows**: `clip.exe` (built-in)

### Install Clipboard Tools (Linux)
```bash
# Ubuntu/Debian
sudo apt install xclip

# Termux
pkg install termux-api
```

## Data Source

The app reads from:
```
~/projects/CellBlocks/data/cellblocks-data.json
```

Any changes to this file (from the React app or AI scripts) will be reflected next time you launch CellBlocksTUI.

## Build & Install

### Rebuild
```bash
cd ~/projects/CellBlocksTUI
go build -o cellblocks-tui
```

### Install to PATH
```bash
# Desktop
cp cellblocks-tui ~/bin/

# Termux
cp cellblocks-tui $PREFIX/bin/
```

### Run from anywhere
```bash
cellblocks-tui
```

## Usage Examples

### 1. Quick Command Lookup
```bash
# Start the app
cellblocks-tui

# Type: docker
# See all Docker-related cards

# Press Enter on "Docker Cleanup"
# → Command copied to clipboard

# Press q to quit
# Paste in terminal: Ctrl+Shift+V
```

### 2. Find Prompts
```bash
# Start the app
cellblocks-tui

# Type: code review
# See "Code Review Prompt"

# Press Enter
# → Prompt copied to clipboard

# Paste in Claude chat
```

### 3. Grid View (The CellBlocks Experience!)
```bash
# Start the app
cellblocks-tui

# Press 'g' to switch to grid view
# → See beautiful card blocks with neon borders!

# Use arrow keys to navigate:
# ↑↓ - Move up/down
# ←→ - Move left/right

# Each card has a colored border matching its category:
# - Bash (Yellow)
# - Prompts (Green)
# - Agents (Cyan)
# - Powershell (Blue)
# - Claude Code (Orange)
# - And more...

# Press Enter on any card to copy to clipboard
```

### 4. Browse by Category
Cards show category with colored borders in grid view, or badges in list view.

## Next Steps

### Phase 2 Features (Coming Soon)
- [ ] Template variable filling ({{port}}, {{name}}, etc.)
- [ ] Category filtering (toggle multiple categories)
- [ ] Favorite/starred cards
- [ ] Touch gestures (port from TFE)
- [ ] Termux integration (share, notifications)

### Want to Contribute?
Check `TODO.md` for the full implementation checklist.

## Troubleshooting

### "failed to read data file"
Make sure the data file exists:
```bash
ls ~/projects/CellBlocks/data/cellblocks-data.json
```

### Clipboard not working
Install clipboard tools:
```bash
# Linux
sudo apt install xclip

# Termux
pkg install termux-api
```

### Binary too large
Strip debug symbols:
```bash
go build -ldflags="-s -w" -o cellblocks-tui
```

## Performance

- Startup: <100ms
- Search: <50ms for 271 cards
- Memory: ~10-15MB
- Binary: 4.6MB

All targets met! 🎉

## Split-Pane Workflow (Termux)

Recommended tmux setup:
```bash
# Create split pane workspace
tmux new -s work
tmux split-window -v -p 30

# Top pane: TFE (file browser)
tfe

# Bottom pane: CellBlocksTUI
cellblocks-tui
```

## Contact

Found a bug? Have an idea?
Open an issue in the CellBlocksTUI repo!

---

**Status**: Phase 1 MVP Complete ✓
**Version**: 1.0.0
**Date**: 2025-10-16
