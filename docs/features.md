# PI-mc — Features

## Left Pane: Dashboard (35% width, 3 stacked sections)

### Queue Section

Shows pending messages queued for the agent.

| Feature | Status | Details |
|---|---|---|
| Show steering messages | Planned | From `queue_update.steering[]` — delivered after current turn |
| Show follow-up messages | Planned | From `queue_update.followUp[]` — delivered when agent finishes |
| Live update on queue_change | Planned | Reacts to every `queue_update` RPC event |
| Visual priority | Planned | Steering shown above follow-up, with distinct icons (`↻` vs `⏳`) |
| Queue count in title | Planned | Title: `Queue (3)` |
| Clear all queued | Future | Abort all pending messages |
| Reorder queue items | Future | Drag with keyboard (Ctrl+↑/↓) |
| Inject new steering message | Future | Type directly into queue and send |

### Steps Section

Shows work steps declared by the agent via `declare_steps`.

| Feature | Status | Details |
|---|---|---|
| Show declared steps | Planned | Parsed from `declare_steps` tool call |
| Live status icons | Planned | `[✓]` done, `[·]` running, `[✗]` failed, `[ ]` pending |
| Auto-advance on tool events | Planned | Steps advance when tool category changes |
| Tool count per step | Planned | `3t` = 3 tool calls in this step |
| Duration per step | Planned | `12s` elapsed |
| Color coding | Planned | Green = done, red = failed, cyan = running |
| Collapse completed steps | Future | Only show current + upcoming |
| Manual step navigation | Future | Jump to specific step |
| Step details on Enter | Future | Show tool calls + output for selected step |
| Re-run step | Future | Ask agent to re-execute a failed step |

### Sessions Section

Lists sessions with navigation, fork, and switch.

| Feature | Status | Details |
|---|---|---|
| List recent sessions | Planned | From `SessionManager.list()` (SDK) or session directory |
| Navigate with ↑↓ | Planned | Keyboard navigation |
| Switch session on Enter | Planned | Sends `switch_session` RPC command |
| Active session marker | Planned | `·` next to current session |
| Session metadata | Planned | Name, message count, last active time |
| Fork from highlighted | Planned | F-key or keybinding to fork selected session |
| New session | Planned | F-key or keybinding |
| Search/filter sessions | Future | Type to filter by name |
| Delete session | Future | Remove old sessions |
| Sort by name/date/activity | Future | Toggle sort mode |

## Right Pane: Chat (65% width)

### Message Display

| Feature | Status | Details |
|---|---|---|
| Streaming text from agent | Planned | Handles `text_delta` events in real-time |
| Thinking block display | Planned | Collapsed by default, expandable |
| Tool call/results inline | Planned | Shows tool name, args summary, result |
| User message display | Planned | `▸` prefix, distinct style |
| Assistant message display | Planned | `●` prefix, color-coded |
| Error message display | Planned | `✗` prefix, red |
| Markdown rendering | Planned | Bold, italic, code blocks, links |
| Syntax highlighting in code blocks | Future | Language-aware coloring |
| Image display | Future | Kitty/iTerm2 inline image protocol |
| Message timestamp toggle | Future | Show/hide per-message timestamps |

### Scrollback

| Feature | Status | Details |
|---|---|---|
| Auto-scroll to bottom | Planned | Scrolls on new content while at bottom |
| Pause auto-scroll | Planned | Stops when user scrolls up manually |
| PgUp/PgDn | Planned | Page-scroll through history |
| Home/End | Planned | Jump to top/bottom |
| Scroll indicator | Planned | `↑ more above` / `↓ more below` |
| Infinite scrollback | Planned | Ring buffer, configurable max lines |
| Jump to message | Future | `:` then line number |
| Search in history | Future | `/` to search, `n`/`N` for next/prev |
| Filter by message type | Future | Show only user messages, only errors, etc. |

### Text Input

| Feature | Status | Details |
|---|---|---|
| Multi-line input | Planned | Shift+Enter for new line, Enter to send |
| Emacs-style editing | Planned | Ctrl+A/E, Ctrl+K/U, Ctrl+W, Alt+F/B |
| Vim-style editing | Future | Modal editing via custom editor |
| Input history | Future | ↑↓ to recall previous prompts |
| Paste support | Planned | Standard terminal paste |
| External editor (Ctrl+G) | Planned | Opens `$VISUAL` or `$EDITOR` |
| Character counter | Future | Shows remaining tokens/characters |
| Command mode (`/`) | Future | Type `/model`, `/compact`, etc. directly |

### Progress Display (inline in chat)

| Feature | Status | Details |
|---|---|---|
| Step status in chat | Planned | `[✓]` `[·]` `[✗]` inline with messages |
| Tool execution spinner | Planned | Spinner while tool is running |
| Duration display | Planned | `3.2s` after tool completes |
| Collapse tool output | Planned | Show summary, expand on Enter |
| Tool output truncation | Planned | `… (full output in /tmp/pi-xxx.log)` |

## Footer: F-Key Bar

| F-Key | Label | Action |
|---|---|---|
| F1 | Queue | Focus Queue pane |
| F2 | Model | Cycle model / open model selector |
| F3 | Steps | Focus Steps pane |
| F4 | Sessions | Focus Sessions pane |
| F5 | Files | Toggle file tree pane (future) |
| F6 | Tree | Open conversation tree navigator (future) |
| F7 | Vim | Open current file in `$EDITOR` (future) |
| F8 | Compact | Trigger manual compaction |
| F9 | Tasks | Open kanban/task board (future) |
| F10 | Quit | Exit PI-mc |

## Global Keybindings

| Key | Action |
|---|---|
| Tab | Cycle focus between panes |
| Escape | Return focus to Chat input |
| Ctrl+C | Abort agent / clear input |
| Ctrl+D | Exit (when input empty) |
| Ctrl+G | Open external editor |
| Ctrl+L | Clear screen / refresh |
| Ctrl+P | Cycle model |
| Shift+Tab | Cycle thinking level |

## Future Panes

These are out of scope for initial release but planned:

| Pane | Description |
|---|---|
| File Tree | Navigable project explorer (like MC's directory browser) |
| Conversation Tree | Visual session tree with branching |
| Tool Output Log | Full tool output viewer with search |
| Resource Monitor | Token usage, cost, context window percentage |
| Task Board | Kanban-style task tracking |
| Diff Viewer | Side-by-side git diff within the TUI |
