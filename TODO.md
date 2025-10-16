# CellBlocksTUI - Implementation Checklist

## Phase 1: Core Functionality (MVP)

### Setup & Foundation
- [x] Create project structure
- [x] Write PLAN.md
- [x] Write README.md
- [x] Initialize go.mod dependencies (`go mod tidy`)
- [x] Create main.go (Bubbletea entry point)
- [x] Create types.go (Card, Category, AppState)
- [x] Create model.go (initialization)

### Data Layer
- [x] Create storage.go
  - [x] LoadData() - read cellblocks-data.json
  - [x] SaveData() - write cellblocks-data.json (future)
  - [x] Expand ~ to home directory
  - [x] Handle missing file gracefully
  - [x] Parse JSON correctly

### Search & Filter
- [x] Create search.go
  - [x] SearchCards() - full-text search
  - [x] FilterByCategories() - category filtering
  - [x] Case-insensitive matching
  - [x] Search title + content

### UI Components
- [x] Create view.go
  - [x] renderListView() - main card list
  - [x] renderGridView() - grid card layout
  - [x] renderPreviewPane() - selected card preview
  - [x] renderHeader() - search + filters
  - [x] renderStatusBar() - help hints
  - [x] renderHelp() - help dialog

- [x] Create styles.go
  - [x] Define color palette
  - [x] Card styles (list and grid)
  - [x] Header styles
  - [x] Status bar styles
  - [x] Grid card styles with category colors

### Event Handling
- [x] Create update.go
  - [x] Window resize
  - [x] Message dispatcher
  - [x] Init() function
  - [x] Keyboard navigation (↑↓←→, jk, hl)
  - [x] Search (/)
  - [x] Copy (Enter)
  - [x] Quit (q)
  - [x] Help (?)
  - [x] Toggle preview (p)
  - [x] Toggle grid view (g)
  - [x] Grid navigation (2D movement)

- [x] Enhanced mouse/touch (ported from TFE)
  - [x] Basic mouse wheel scrolling
  - [x] Click to select card (list and grid view)
  - [x] Double-click to copy (500ms threshold)
  - [x] Grid click detection (X/Y coordinate calculation)
  - [x] WithMouseAllMotion() for better tracking
  - [ ] Swipe gestures (future enhancement)

### Integration
- [x] Create clipboard.go
  - [x] Detect platform (Termux/Linux/macOS)
  - [x] CopyToClipboard() - termux-clipboard-set
  - [x] Fallback to xclip/pbcopy
  - [x] Windows support (clip.exe)

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

**Phase 1 Progress:** ✅ ~95% Complete!
**Completed:**
- All core functionality (data loading, search, filtering)
- Full UI implementation (list view, grid view, preview panes)
- Complete keyboard navigation
- Enhanced mouse/touch navigation (click, double-click, wheel scroll)
- Clipboard integration (all platforms)
- Grid view selection visibility fix

**Remaining:**
- Testing with real cellblocks-data.json
- Swipe gestures (optional enhancement)

**Next Tasks:**
1. Test with actual cellblocks-data.json (validate with 234+ cards)
2. Begin Phase 2 (template support - {{variable}} detection)
3. Add swipe gestures (optional)

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
