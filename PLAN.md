# CellBlocksTUI - Terminal Card Management System

**Version**: 1.1.0
**Status**: Production Ready (Phase 2 Complete!)
**Author**: Matt
**Created**: 2025-10-16
**Last Updated**: 2025-10-17

## Overview

CellBlocksTUI is a terminal-based interface for managing and accessing CellBlocks cards - your library of commands, prompts, agent configurations, and code snippets. Built specifically for **Termux mobile workflows** and **touch-friendly navigation**.

### Key Goals

1. **Lightweight** - 5MB binary, ~10MB RAM (vs 110MB React version)
2. **Offline-first** - Works without server, reads local JSON
3. **Touch-optimized** - Port TFE's proven touch navigation
4. **Termux-native** - Clipboard, share, notifications via termux-api
5. **Split-pane ready** - Perfect companion to TFE in tmux
6. **Data compatible** - Shares cellblocks-data.json with React version

## Architecture

### Data Flow

```
┌─────────────────────────────────────────────────┐
│  Desktop (React CellBlocks)                     │
│  - Create cards via AI scripts                  │
│  - Rich editing in browser                      │
│  - Socket.io real-time updates                  │
└─────────────┬───────────────────────────────────┘
              │
              │ Syncs via Tailscale/Syncthing
              ▼
┌─────────────────────────────────────────────────┐
│  Mobile (CellBlocksTUI)                         │
│  - Quick lookups in Termux                      │
│  - Touch navigation (TFE patterns)              │
│  - Offline access                               │
│  - Clipboard integration                        │
└─────────────────────────────────────────────────┘

Both share:
~/projects/CellBlocks/data/cellblocks-data.json
```

### File Structure

```
CellBlocksTUI/
├── main.go                 - Entry point (21 lines, Bubbletea setup)
├── types.go                - Card, Category, AppState structs
├── model.go                - State initialization & layout
├── update.go               - Message dispatcher
├── update_keyboard.go      - Keyboard shortcuts
├── update_mouse.go         - Touch/mouse navigation (from TFE)
├── view.go                 - Layout & rendering
├── storage.go              - Read/write cellblocks-data.json
├── search.go               - Full-text search engine
├── clipboard.go            - Termux clipboard integration
├── socket.go               - (Optional) Socket.io live sync
├── styles.go               - Lipgloss theme
├── go.mod                  - Dependencies
├── PLAN.md                 - This file
└── README.md               - User documentation

Future:
├── edit.go                 - Card creation/editing
└── sync.go                 - Syncthing integration helpers
```

## Data Model

### JSON Structure (Shared with React)

```json
{
  "cards": [
    {
      "id": "card-1234567890",
      "title": "Docker Run Command",
      "content": "docker run -p {{port}}:{{port}} --name {{container}} {{image}}",
      "categoryId": "cat-bash",
      "createdAt": 1697234567890,
      "updatedAt": 1697234567890,
      "imageId": null
    }
  ],
  "categories": [
    {
      "id": "cat-bash",
      "name": "Bash",
      "color": "#ffff00",
      "hidden": false,
      "parentCategoryId": null
    },
    {
      "id": "cat-prompts",
      "name": "Prompts",
      "color": "#00ff41",
      "hidden": false,
      "parentCategoryId": null
    }
  ]
}
```

### Go Structs

```go
type Card struct {
    ID         string `json:"id"`
    Title      string `json:"title"`
    Content    string `json:"content"`
    CategoryID string `json:"categoryId"`
    CreatedAt  int64  `json:"createdAt"`
    UpdatedAt  int64  `json:"updatedAt"`
    ImageID    string `json:"imageId,omitempty"`
}

type Category struct {
    ID               string `json:"id"`
    Name             string `json:"name"`
    Color            string `json:"color"`
    Hidden           bool   `json:"hidden,omitempty"`
    ParentCategoryID string `json:"parentCategoryId,omitempty"`
}

type CellBlocksData struct {
    Cards      []Card     `json:"cards"`
    Categories []Category `json:"categories"`
}

type AppState struct {
    // Data
    Data           *CellBlocksData
    FilteredCards  []Card

    // UI State
    SelectedIndex  int
    ScrollOffset   int
    SearchQuery    string
    SelectedCategories []string

    // View mode
    ViewMode       ViewMode  // list, grid, detail
    ShowPreview    bool

    // Template editing
    TemplateVars   map[string]string

    // Terminal size
    Width          int
    Height         int
}

type ViewMode int
const (
    ViewList ViewMode = iota
    ViewGrid
    ViewDetail
)
```

## Features

### Phase 1: Core Functionality ✅ COMPLETE!

**Goal**: Read-only card browsing with search

- [x] Project setup from TUITemplate
- [x] Load cellblocks-data.json (234+ cards tested)
- [x] Display cards in list and grid views
- [x] Category filtering (toggle multiple)
- [x] Full-text search (title + content)
- [x] Touch navigation (ported from TFE)
- [x] Preview pane for selected card with scrolling
- [x] Copy card to clipboard
- [x] Termux clipboard integration (all platforms)
- [x] Help dialog (? key)
- [x] Responsive layout (mobile-optimized)
- [x] Mouse wheel scrolling (position-aware)
- [x] Grid view with category-colored borders

**Deliverable**: ✅ Functional card browser for Termux

### Phase 2: Enhanced Features ✅ COMPLETE!

**Goal**: Category filtering UI and card creation

- [x] Interactive category filter screen (f key)
  - [x] Checkbox list with navigation
  - [x] Toggle multiple categories
  - [x] Select all / Clear all
  - [x] Mobile-friendly header display
- [x] Card creation form (n key)
  - [x] Multi-field form (Title, Content, Category)
  - [x] Tab navigation between fields
  - [x] Real-time validation
  - [x] Save to JSON file
  - [x] Auto-jump to new card
- [x] Auto-reload functionality
  - [x] Check for file changes every 10 seconds
  - [x] Visual notification for new cards
  - [x] Perfect for AI-generated cards
- [x] Bug fixes
  - [x] Fixed flickering filter text
  - [x] Fixed footer duplication on narrow screens
  - [x] Fixed mouse wheel scroll locking
  - [x] Responsive status bar (3 breakpoints)

**Deliverable**: ✅ Full CRUD operations + auto-sync

### Phase 3: Template Support ✅ COMPLETE!

**Goal**: Fill in template variables

- [x] Detect {{variable}} syntax
- [x] Show input fields for variables
- [x] Real-time template preview
- [x] Copy filled template to clipboard
- [x] Save filled values for reuse
- [x] Full-screen detail view with markdown rendering
- [x] Support for default values {{var|default}}

**Deliverable**: ✅ Interactive template filling with ViewDetail mode!

### Phase 4: Enhanced Termux Integration (Future)

**Goal**: Native mobile experience

- [x] termux-clipboard-set integration ✅
- [ ] termux-share (share to other apps)
- [ ] termux-notification (card updates)
- [ ] termux-toast (quick messages)
- [ ] Launch URLs via termux-open-url
- [ ] Execute commands directly

**Deliverable**: Full Termux native experience

### Phase 5: Optional Features

- [ ] Socket.io listener (live updates from React)
- [x] File watcher (auto-reload) ✅
- [x] Card creation ✅
- [ ] Card editing
- [ ] Card deletion
- [ ] Favorites/starred cards
- [ ] Recent cards history
- [ ] Export search results
- [ ] Syncthing integration helper script

## UI Design

### Mobile Layout (Termux)

```
┌─ CellBlocks TUI ─────────────────────────┐
│ 🔍 Search: docker▊               [234]   │  ← Compact header
├──────────────────────────────────────────┤
│ ☰ Categories                      [3/8]  │  ← Tap to expand
│ ✓ Bash  ✓ Prompts  ✓ Agents             │  ← Toggle filters
├──────────────────────────────────────────┤
│ Docker Commands                    Bash  │  ← Card list
│ > docker run -p {{port}}:{{port}}...     │     (scrollable)
│ ─────────────────────────────────────────│
│ Git Workflow                    Prompts  │
│   Step-by-step git PR workflow...        │
│ ─────────────────────────────────────────│
│ Code Review Checklist            Agents  │
│   Review for security, performance...    │
├──────────────────────────────────────────┤
│ Preview: Docker Commands            [📋] │  ← Preview pane
│ docker run -p {{port}}:{{port}} \        │     (collapsible)
│   --name {{container}} {{image}}         │
│                                          │
│ Variables:                               │
│ port: 3000▊                              │  ← Input fields
│ container: my-app                        │
│ image: nginx:latest                      │
│                                          │
│ → docker run -p 3000:3000 \              │  ← Filled result
│     --name my-app nginx:latest           │
│ [Copy] [Share] [Execute]                 │  ← Actions
└──────────────────────────────────────────┘

? help | n new | / search | c categories | q quit
```

### Touch Gestures (From TFE)

| Gesture | Action |
|---------|--------|
| **Tap** | Select card |
| **Double-tap** | Copy to clipboard |
| **Swipe left/right** | Navigate categories |
| **Swipe up/down** | Scroll card list |
| **Long-press** | Context menu (Edit/Delete/Share) |
| **Pinch** | Zoom text (accessibility) |

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑↓` / `jk` | Navigate cards |
| `←→` / `hl` | Switch categories |
| `/` | Search |
| `Enter` | Copy selected card |
| `Space` | Toggle category filter |
| `p` | Toggle preview pane |
| `e` | Edit template variables |
| `c` | Copy filled template |
| `s` | Share via Termux |
| `x` | Execute command |
| `?` | Help |
| `q` | Quit |
| `Ctrl+R` | Reload data |
| `Esc` | Clear search/deselect |

## Technical Implementation

### Storage Layer (storage.go)

```go
const DataPath = "~/projects/CellBlocks/data/cellblocks-data.json"

func LoadData() (*CellBlocksData, error)
func SaveData(data *CellBlocksData) error
func WatchFile(callback func()) error  // Optional: auto-reload
```

### Search Engine (search.go)

```go
func SearchCards(cards []Card, query string) []Card
func FilterByCategories(cards []Card, categoryIDs []string) []Card
func SortCards(cards []Card, sortBy SortOption) []Card

// Full-text search algorithm (from CellBlocks React)
// - Lowercase matching
// - Search both title and content
// - Highlight matches in results
```

### Clipboard Integration (clipboard.go)

```go
func CopyToClipboard(text string) error {
    // Termux: termux-clipboard-set
    // Linux: xclip, xsel
    // macOS: pbcopy
    // Windows: clip.exe
}

func ShareToApp(text string) error {
    // Termux only: termux-share
}

func ShowNotification(title, text string) error {
    // Termux: termux-notification
}
```

### Template Parser (template.go)

```go
func ExtractVariables(content string) []string {
    // Find all {{variable}} patterns
    // Return unique variable names
}

func FillTemplate(content string, vars map[string]string) string {
    // Replace {{variable}} with values
    // Support default values: {{port|3000}}
}

func ParseDefaultValue(variable string) (name, defaultVal string) {
    // Split "port|3000" -> ("port", "3000")
}
```

### Socket.io Client (socket.go) - Optional

```go
func ConnectToServer() (*socketio.Client, error)
func ListenForCardUpdates(onUpdate func(Card))
func ListenForCardCreation(onCreate func(Card))

// Only needed if you want live sync from React version
// Can be disabled for pure offline mode
```

## Dependencies

```go
// go.mod
module github.com/Matt/cellblocks-tui

go 1.22

require (
    github.com/charmbracelet/bubbletea v0.27.0
    github.com/charmbracelet/lipgloss v0.13.1
    github.com/charmbracelet/bubbles v0.20.0

    // Optional:
    // github.com/googollee/go-socket.io v1.7.0  // For live sync
    // github.com/fsnotify/fsnotify v1.7.0        // For file watching
)
```

## Development Roadmap

### Week 1: Foundation
- [x] Set up project from TUITemplate
- [ ] Define data structures
- [ ] Implement data loading
- [ ] Basic list view
- [ ] Card selection

### Week 2: Core Features
- [ ] Full-text search
- [ ] Category filtering
- [ ] Preview pane
- [ ] Touch navigation (port from TFE)
- [ ] Clipboard integration

### Week 3: Templates & Polish
- [ ] Template variable extraction
- [ ] Interactive variable input
- [ ] Template filling
- [ ] Help system
- [ ] Error handling

### Week 4: Termux Integration
- [ ] termux-api integration
- [ ] Share functionality
- [ ] Notifications
- [ ] Mobile layout optimization
- [ ] Testing on actual device

## Testing Strategy

### Manual Testing Checklist

**Data Loading:**
- [ ] Loads cellblocks-data.json correctly
- [ ] Handles missing file gracefully
- [ ] Parses all 234+ cards from sample data
- [ ] Displays categories correctly

**Search:**
- [ ] Case-insensitive search works
- [ ] Searches both title and content
- [ ] Updates results in real-time
- [ ] Shows result count

**Navigation:**
- [ ] Arrow keys move cursor
- [ ] vim keys (hjkl) work
- [ ] Touch taps select cards
- [ ] Scroll works smoothly

**Templates:**
- [ ] Detects {{variable}} syntax
- [ ] Shows input fields
- [ ] Fills template correctly
- [ ] Copies filled result

**Termux:**
- [ ] Clipboard integration works
- [ ] Share to apps works
- [ ] Notifications appear
- [ ] Touch gestures responsive

### Device Testing

- [ ] Termux on Android (phone)
- [ ] Termux on Android (tablet)
- [ ] Linux desktop (fallback)
- [ ] macOS (if available)

## Performance Goals

- **Startup time**: < 100ms
- **Search latency**: < 50ms for 1000 cards
- **Memory usage**: < 15MB
- **Binary size**: < 5MB
- **Render time**: < 16ms (60 FPS)

## Data Sync Strategy

### Option 1: Syncthing (Recommended)
```bash
# Install Syncthing on both devices
pkg install syncthing  # Termux
brew install syncthing # Desktop

# Sync ~/projects/CellBlocks/data/ folder
# Auto-syncs changes in both directions
# Works offline (syncs when connected)
```

### Option 2: Rsync + Tailscale
```bash
# Manual sync when connected
alias sync-cellblocks='rsync -avz \
  desktop:~/projects/CellBlocks/data/ \
  ~/CellBlocks/data/'
```

### Option 3: Git
```bash
# Treat data as git repo
cd ~/projects/CellBlocks/data
git init
git add cellblocks-data.json
git commit -m "Update cards"
git push

# Pull on mobile
cd ~/CellBlocks/data
git pull
```

## Deployment

### Build
```bash
cd ~/projects/CellBlocksTUI
go build -o cellblocks-tui
```

### Install (Desktop)
```bash
cp cellblocks-tui ~/bin/
chmod +x ~/bin/cellblocks-tui
```

### Install (Termux)
```bash
cp cellblocks-tui $PREFIX/bin/
chmod +x $PREFIX/bin/cellblocks-tui
```

### Launch Script (tmux workspace)
```bash
# Save to ~/bin/workspace
#!/bin/bash
tmux new -s cellblocks \; \
  split-window -v -p 30 \; \
  send-keys -t 0 'tfe' C-m \; \
  send-keys -t 1 'cellblocks-tui' C-m \; \
  select-pane -t 0
```

## Success Metrics

- [ ] Successfully reads 234+ cards from cellblocks-data.json
- [ ] Search returns results in < 50ms
- [ ] Touch navigation feels native (TFE quality)
- [ ] Copy to clipboard works 100% of time
- [ ] Memory usage stays under 15MB
- [ ] Works offline without degradation
- [ ] Launches in < 100ms

## Future Enhancements

### Post-MVP Features
- Card editing (create/update/delete)
- Sync status indicator
- Favorites/starred cards
- Recently viewed cards
- Card versioning
- Export filtered results
- Import from clipboard
- Voice input (termux-speech-to-text)

### Advanced Features
- AI card creation from TUI
- Custom themes
- Plugin system
- Card statistics
- Backup/restore
- Multi-device sync indicator
- Conflict resolution for edits

## Known Limitations

1. **Read-only in MVP** - No editing initially (React for that)
2. **No images** - Text-only cards (images stay in React version)
3. **Single user** - No collaboration features
4. **No version history** - Overwrites data file
5. **Limited formatting** - Plain text rendering only

## References

### Existing Projects
- **CellBlocks React**: ~/projects/CellBlocks/
- **TFE**: ~/projects/TFE/ (touch navigation patterns)
- **CSV Viewer**: ~/projects/csv-viewer/ (grid layout ideas)
- **TUITemplate**: ~/projects/TUITemplate/ (architecture base)

### Documentation
- CellBlocks CLAUDE.md - AI integration details
- TFE CLAUDE.md - Touch navigation implementation
- Bubbletea docs - https://github.com/charmbracelet/bubbletea
- Lipgloss docs - https://github.com/charmbracelet/lipgloss

### Data Files
- Data source: ~/projects/CellBlocks/data/cellblocks-data.json
- AI scripts: ~/projects/CellBlocks/scripts/ai-*.js
- Socket.io server: ~/projects/CellBlocks/server/pty-server.ts

---

**Status**: ✅ Phase 3 Complete - Template Support Ready!

**Completed**:
- ✅ Phase 1: Core Functionality (100%)
- ✅ Phase 2: Enhanced Features (100%)
  - Category filtering UI
  - Card creation
  - Auto-reload
  - Mobile optimizations
  - Bug fixes
- ✅ Phase 3: Template Support (100%)
  - Full-screen ViewDetail mode
  - Glamour markdown rendering
  - Template variable detection and filling
  - Interactive template form
  - Real-time preview

**What's Working**:
- 234+ cards tested from real cellblocks-data.json
- Full-text search with <50ms latency
- Category filtering with interactive UI
- Card creation with validation
- Auto-reload every 10 seconds (perfect for AI cards!)
- Mouse & keyboard navigation
- Preview pane with position-aware scrolling
- Multi-platform clipboard (Termux/Linux/macOS/Windows)
- Mobile-responsive design (60-120+ char widths)
- Grid and list views with adaptive layouts
- **NEW**: Full-screen detail view (Enter or 'd' key)
- **NEW**: Markdown rendering with Glamour (toggle with 'm')
- **NEW**: Template variable detection ({{variable}} and {{var|default}})
- **NEW**: Interactive template form with Tab navigation
- **NEW**: Real-time filled template preview
- **NEW**: Copy filled templates to clipboard

**Next Steps**:
1. **Phase 4: Enhanced Termux** - Add share, notifications, toasts
2. **Card Editing** - Edit existing cards from TUI
3. **Favorites System** - Star frequently used cards
4. **Card Deletion** - Delete cards from TUI

**Performance Metrics**:
- Binary size: ~5MB
- Memory usage: ~10-12MB
- Startup time: <100ms
- Render time: <16ms (60 FPS)

**Questions/Blockers**: None!
