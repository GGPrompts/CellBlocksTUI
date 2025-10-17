# CellBlocksTUI

A lightweight terminal interface for managing CellBlocks cards - your command library, prompts, and code snippets.

**Perfect for Termux mobile workflows!** 🚀📱

## Features

✅ **Lightweight** - 5MB binary, ~10MB RAM (vs 110MB React version)
✅ **Touch-optimized** - Enhanced mouse/touch navigation with click, double-click, and wheel scrolling
✅ **Offline-first** - No server required, works completely offline
✅ **Auto-reload** - Detects new cards every 10 seconds (perfect for AI-generated cards!)
✅ **Category filtering** - Interactive UI to filter by multiple categories
✅ **Card creation** - Create new cards directly from the TUI
✅ **Termux-native** - Clipboard integration (share & notifications coming soon!)
✅ **Data compatible** - Shares `cellblocks-data.json` with React CellBlocks
✅ **Split-pane ready** - Perfect companion to TFE in tmux
✅ **Mobile-responsive** - Adapts UI to screen width (tested down to 60 chars)

## Quick Start

```bash
# From anywhere in Termux or desktop:
cellblocks-tui

# Or with TFE in split pane:
~/bin/workspace  # Launches TFE + CellBlocksTUI
```

## Installation

### Desktop
```bash
cd ~/projects/CellBlocksTUI
go build -o cellblocks-tui
cp cellblocks-tui ~/bin/
```

### Termux
```bash
cd ~/projects/CellBlocksTUI
go build -o cellblocks-tui
cp cellblocks-tui $PREFIX/bin/
```

## Usage

### Keyboard Shortcuts

**Navigation:**
- `↑↓` or `j/k` - Navigate cards
- `←→` or `h/l` - Navigate categories (grid view)
- `Shift+↑↓` - Scroll preview content
- `PgUp/PgDn` - Scroll by page
- `Home/End` - Jump to first/last card

**View:**
- `g` - Toggle grid/list view
- `p` - Toggle preview pane (side-by-side on wide screens!)
- `Space` - Pin card to preview (grid view)

**Actions:**
- `Enter` or `c` - Copy card to clipboard
- `n` - Create new card
- `f` - Filter by category
- `/` - Clear filters
- `Type...` - Search cards (real-time)
- `Backspace` - Delete search character

**General:**
- `?` - Show help
- `Esc` - Close help/clear search/exit screens
- `q` or `Ctrl+C` - Quit

### Mouse & Touch

**Mouse:**
- **Click** - Select card & pin to preview
- **Double-click** - Copy card to clipboard
- **Wheel scroll over list** - Scroll cards (preview stays locked)
- **Wheel scroll over preview** - Scroll preview content

**Touch (Termux):**
- **Tap** - Select card
- **Double-tap** - Copy to clipboard
- **Swipe** - Scroll lists and preview

## New Features

### Category Filtering (Press `f`)
- Interactive checkbox interface
- Toggle multiple categories with `Space/Enter`
- `a` - Select all categories
- `c` - Clear all filters
- Shows active filter count in header
- Mobile-friendly display (shows count on narrow screens)

### Card Creation (Press `n`)
- Multi-field form: Title, Content, Category
- Tab through fields with `Tab/Shift+Tab`
- Multi-line content support (press `Enter` for newlines)
- Select category with `↑↓` arrows
- Real-time validation
- Save with `Ctrl+S` or `Ctrl+Enter`
- Automatically jumps to new card after save

### Auto-Reload
- Checks for file changes every 10 seconds
- Shows notification when new cards detected: "✨ 3 new card(s) detected!"
- Perfect for monitoring AI-generated cards
- Notification auto-dismisses after 5 seconds
- No need to restart TUI when React app or AI adds cards!

## Data Source

Reads from: `~/projects/CellBlocks/data/cellblocks-data.json`

**Important:** This is the same file used by the React version! Both can run simultaneously, and changes sync automatically via the auto-reload feature.

## Syncing Data

### With Syncthing (Recommended)
```bash
# Desktop + Mobile auto-sync
pkg install syncthing  # Termux
# Sync ~/projects/CellBlocks/data/
```

### With Rsync
```bash
# Manual sync via Tailscale
rsync -avz desktop:~/projects/CellBlocks/data/ ~/CellBlocks/data/
```

## Split Pane Workspace

Create `~/bin/workspace`:
```bash
#!/bin/bash
tmux new -s work \; \
  split-window -v -p 30 \; \
  send-keys -t 0 'tfe' C-m \; \
  send-keys -t 1 'cellblocks-tui' C-m \; \
  select-pane -t 0
```

Then just run: `workspace`

## Mobile Optimization

The UI adapts to screen width:

- **< 60 chars** (phone portrait): Shows only essential hints
- **60-90 chars** (phone landscape): Shows most common actions
- **> 90 chars** (desktop/tablet): Shows all hints and features
- **> 120 chars** (wide desktop): Preview pane appears side-by-side with grid view

Header filters also adapt:
- Narrow screens: "Filters: 3"
- Wide screens: "Filters: AWS, Docker +1"

## Architecture

Built with **Bubbletea** (Elm architecture for Go):

```
CellBlocksTUI/
├── main.go              - Entry point
├── types.go             - Data structures
├── model.go             - State management
├── update.go            - Event handling
├── update_mouse.go      - Mouse/touch navigation
├── view.go              - Rendering
├── storage.go           - File I/O & auto-reload
├── search.go            - Search & filtering
├── clipboard.go         - Multi-platform clipboard
└── styles.go            - Lipgloss theming
```

See [PLAN.md](./PLAN.md) for complete technical details.

## Development Status

### ✅ Phase 1: Core Functionality (COMPLETE!)
- [x] Load cellblocks-data.json (234+ cards tested)
- [x] Display cards in list and grid views
- [x] Full-text search (title + content)
- [x] Touch/mouse navigation (click, double-click, wheel)
- [x] Clipboard integration (Termux/Linux/macOS/Windows)
- [x] Preview pane with scrolling (adaptive layouts)
- [x] Help system

### ✅ Phase 2: Enhanced Features (COMPLETE!)
- [x] Category filtering with interactive UI
- [x] Card creation with form validation
- [x] Auto-reload for external changes
- [x] Mobile-responsive design
- [x] Bug fixes (flickering, scroll locking, footer duplication)

### 🔜 Phase 3: Template Support (Next!)
- [ ] Detect `{{variable}}` syntax
- [ ] Interactive variable input fields
- [ ] Real-time template preview
- [ ] Copy filled templates

### 🔜 Phase 4: Enhanced Termux Integration
- [x] termux-clipboard-set (DONE!)
- [ ] termux-share (share to other apps)
- [ ] termux-notification (card updates)
- [ ] termux-toast (quick messages)
- [ ] termux-open-url (open links)

## Performance

- **Binary size:** ~5MB
- **Memory usage:** ~10-12MB
- **Startup time:** <100ms
- **Search latency:** <50ms for 234 cards
- **Render time:** <16ms (60 FPS)
- **Auto-reload:** Checks every 10 seconds

## Recent Updates

### v1.1.0 (Latest)
- ✨ Added category filtering UI (press `f`)
- ✨ Added card creation form (press `n`)
- ✨ Added auto-reload every 10 seconds
- 🐛 Fixed flickering filter text (alphabetical sorting)
- 🐛 Fixed footer duplication on narrow screens
- 🐛 Fixed mouse wheel scrolling both list and preview
- 📱 Made status bar responsive (3 width breakpoints)
- 🎨 Improved header to show filter count on mobile

### v1.0.0 (Initial Release)
- 🎉 Full card browsing with grid and list views
- 🔍 Real-time search across 234+ cards
- 🖱️ Enhanced mouse/touch navigation
- 📋 Multi-platform clipboard support
- 📱 Mobile-optimized for Termux

## Troubleshooting

### Binary is busy when copying
```bash
# Kill the running instance first
pkill cellblocks-tui

# Then copy
cp cellblocks-tui ~/bin/cellblocks-tui
```

### Clipboard not working on Termux
```bash
# Install termux-api
pkg install termux-api

# Test it
echo "test" | termux-clipboard-set
termux-clipboard-get
```

### File not loading
Check that the data file exists:
```bash
ls -lh ~/projects/CellBlocks/data/cellblocks-data.json
```

## Contributing

This is a personal project, but feel free to fork and adapt for your needs!

## License

MIT

## Credits

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) by Charm
- Styled with [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Architecture inspired by [TUITemplate](../TUITemplate/)
- Touch patterns ported from [TFE](../TFE/)
- Data compatible with [CellBlocks React](../CellBlocks/)

---

**Made with ❤️ for Termux mobile workflows**
