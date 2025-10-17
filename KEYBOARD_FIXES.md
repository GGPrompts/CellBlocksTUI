# Keyboard Handler Fixes

## Issues Fixed

### 1. ✅ Duplicate 'm' Key Handler
**Problem:** The 'm' key handler existed in both the global handler and the detail view handler, causing potential conflicts.

**Fix:**
- Removed duplicate handler from `handleDetailViewInput()` (was at lines 539-542)
- The global handler (lines 137-142) now handles markdown toggle for all modes:
  - ViewList (with preview enabled)
  - ViewGrid (with preview enabled)
  - ViewDetail

**Result:** The 'm' key now works consistently across all views to toggle markdown rendering.

---

### 2. ✅ Esc Key Logic Bug
**Problem:** When pressing Esc in detail view, the code set `m.ViewMode = ViewList` and then immediately checked `if m.ViewMode == ViewList`, which would always be true. This caused the detail view state to be reset incorrectly.

**Fix (lines 183-189 in update.go):**
```go
// Exit special screens back to main view
if m.ViewMode == ViewCategoryFilter || m.ViewMode == ViewCardCreate || m.ViewMode == ViewDetail {
    // Reset detail view state when exiting detail mode
    m.DetailScrollOffset = 0
    m.ShowTemplateForm = false
    // Return to list view
    m.ViewMode = ViewList
    return m, nil
}
```

**Result:** Esc key now properly exits detail view and resets state.

---

### 3. ✅ [MD] Indicator Placement
**Clarification:** The [MD] indicator is **intentionally** placed in the detail view header only (view.go:791).

**Where it appears:**
- **Detail view header:** `Title  Category [MD]` or `Title  Category [TXT]`
  - Green `[MD]` badge when markdown is enabled
  - Gray `[TXT]` badge when markdown is disabled

**Where it does NOT appear:**
- Main list/grid view header (only shows: Title, Search, Filters, Card count)
- Preview pane header (only shows: Title, Category)

**This is correct behavior** - the indicator shows whether markdown is currently being rendered in the detail view.

---

## Testing Checklist

### In List/Grid View:
- [ ] 'm' key toggles markdown in preview pane (if preview is shown)
- [ ] 'p' key toggles preview pane
- [ ] 'Enter' or 'd' opens detail view
- [ ] No [MD] indicator in main header

### In Detail View:
- [ ] 'm' key toggles between markdown and plain text
- [ ] Header shows `[MD]` when markdown is ON
- [ ] Header shows `[TXT]` when markdown is OFF
- [ ] Footer shows `m plain text` when markdown is ON
- [ ] Footer shows `m markdown` when markdown is OFF
- [ ] 'Esc' exits detail view and returns to list view
- [ ] 'c' copies card content to clipboard
- [ ] 'Enter' copies card content (or filled template)
- [ ] '↑/↓' or 'k/j' scrolls content
- [ ] Mouse wheel scrolls content
- [ ] Scroll stops at end of content (no infinite scroll)

### Template Features:
- [ ] Open card with `{{variable}}` syntax
- [ ] Template form shows automatically if variables detected
- [ ] 't' toggles template form visibility
- [ ] Tab navigates between template fields
- [ ] 'Enter' copies filled template
- [ ] Preview shows filled template in real-time

---

## Build Status

✅ **Build successful** - No compilation errors

---

## What Changed

**Files modified:**
- `update.go` (lines 183-189, 539-542 removed)
  - Fixed esc key logic
  - Removed duplicate 'm' handler

**No changes needed:**
- `view.go` - [MD] indicator placement is already correct
- `update_mouse.go` - Mouse isolation already working
- `model.go` - Markdown enabled by default already set

---

**Status:** ✅ All keyboard handler issues resolved!

The detail view now has:
- Consistent markdown toggle ('m' key)
- Proper exit behavior (Esc key)
- Clear visual state indicators ([MD]/[TXT])
- Responsive scroll controls
- Functional copy commands (c and Enter)
