# CellBlocksTUI

A lightweight terminal interface for managing CellBlocks cards - your command library, prompts, and code snippets.

**Perfect for Termux mobile workflows!** 🚀📱

## Quick Start

```bash
# From anywhere in Termux or desktop:
cellblocks-tui

# Or with TFE in split pane:
~/bin/workspace  # Launches TFE + CellBlocksTUI
```

## Features

✅ **Lightweight** - 5MB binary, ~10MB RAM (vs 110MB React version)
✅ **Touch-optimized** - Proven navigation patterns from TFE
✅ **Offline-first** - No server required
✅ **Termux-native** - Clipboard, share, notifications
✅ **Data compatible** - Shares JSON with React CellBlocks
✅ **Split-pane ready** - Perfect companion to TFE

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

### Navigation
- **↑↓** or **jk** - Move between cards
- **←→** or **hl** - Switch categories
- **/** - Search
- **Enter** - Copy to clipboard
- **Space** - Toggle category filter
- **p** - Toggle preview pane
- **?** - Show help
- **q** - Quit

### Touch Gestures (Termux)
- **Tap** - Select card
- **Double-tap** - Copy to clipboard
- **Swipe** - Scroll/navigate
- **Long-press** - Context menu

## Data Source

Reads from: `~/projects/CellBlocks/data/cellblocks-data.json`

**Important:** This is the same file used by the React version! Both can run simultaneously.

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

## Architecture

See [PLAN.md](./PLAN.md) for complete technical details.

## Development Status

**Phase 1: Core Functionality ✅ Complete!**
- [x] Load cellblocks-data.json (271 cards, 10 categories)
- [x] Display cards in list and grid views
- [x] Full-text search (title + content)
- [x] Category filtering
- [x] Touch/mouse navigation (click, double-click, wheel)
- [x] Clipboard integration (all platforms)
- [x] Preview pane (adaptive layouts)
- [x] Help system

**Ready for testing!** Try running `./cellblocks-tui` or `cellblocks-tui` if installed.

## License

MIT

## Credits

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea)
- Architecture from [TUITemplate](../TUITemplate/)
- Touch patterns from [TFE](../TFE/)
- Data compatible with [CellBlocks React](../CellBlocks/)
