# CellBlocksTUI - Implementation Checklist

## Phase 1: Core Functionality (MVP)

### Setup & Foundation
- [x] Create project structure
- [x] Write PLAN.md
- [x] Write README.md
- [ ] Initialize go.mod dependencies (`go mod tidy`)
- [ ] Create main.go (Bubbletea entry point)
- [ ] Create types.go (Card, Category, AppState)
- [ ] Create model.go (initialization)

### Data Layer
- [ ] Create storage.go
  - [ ] LoadData() - read cellblocks-data.json
  - [ ] SaveData() - write cellblocks-data.json (future)
  - [ ] Expand ~ to home directory
  - [ ] Handle missing file gracefully
  - [ ] Parse JSON correctly

### Search & Filter
- [ ] Create search.go
  - [ ] SearchCards() - full-text search
  - [ ] FilterByCategories() - category filtering
  - [ ] Case-insensitive matching
  - [ ] Search title + content

### UI Components
- [ ] Create view.go
  - [ ] renderListView() - main card list
  - [ ] renderPreviewPane() - selected card preview
  - [ ] renderHeader() - search + filters
  - [ ] renderStatusBar() - help hints

- [ ] Create styles.go
  - [ ] Define color palette
  - [ ] Card styles
  - [ ] Header styles
  - [ ] Status bar styles

### Event Handling
- [ ] Create update.go
  - [ ] Window resize
  - [ ] Message dispatcher
  - [ ] Init() function

- [ ] Create update_keyboard.go
  - [ ] Navigation (↑↓←→, jk, hl)
  - [ ] Search (/)
  - [ ] Copy (Enter)
  - [ ] Quit (q)
  - [ ] Help (?)
  - [ ] Toggle preview (p)

- [ ] Create update_mouse.go (port from TFE)
  - [ ] Click to select
  - [ ] Double-click to copy
  - [ ] Scroll gestures
  - [ ] Touch-friendly detection

### Integration
- [ ] Create clipboard.go
  - [ ] Detect platform (Termux/Linux/macOS)
  - [ ] CopyToClipboard() - termux-clipboard-set
  - [ ] Fallback to xclip/pbcopy

### Testing
- [ ] Test with real cellblocks-data.json
- [ ] Test search functionality
- [ ] Test category filtering
- [ ] Test clipboard on Termux
- [ ] Test touch navigation
- [ ] Test with 234+ cards

## Phase 2: Template Support

- [ ] Create template.go
  - [ ] ExtractVariables() - find {{variable}}
  - [ ] FillTemplate() - replace with values
  - [ ] ParseDefaultValue() - handle {{port|3000}}

- [ ] Add template UI
  - [ ] Input fields for variables
  - [ ] Real-time preview
  - [ ] Copy filled template

## Phase 3: Termux Integration

- [ ] termux-clipboard-set
- [ ] termux-share (share to apps)
- [ ] termux-notification
- [ ] termux-toast
- [ ] termux-open-url

## Phase 4: Optional Features

- [ ] Socket.io listener (live sync)
- [ ] File watcher (auto-reload)
- [ ] Card creation/editing
- [ ] Favorites/starred
- [ ] Recent history
- [ ] Export results

## Current Status

**In Progress:** Project setup and planning
**Next Task:** Create types.go and storage.go
**Blockers:** None

## Quick Start Commands

```bash
# Install dependencies
cd ~/projects/CellBlocksTUI
go mod tidy

# Build
go build -o cellblocks-tui

# Run
./cellblocks-tui

# Install to PATH
cp cellblocks-tui ~/bin/  # Desktop
cp cellblocks-tui $PREFIX/bin/  # Termux
```
