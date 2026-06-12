# PI-mc — Roadmap

## Phase 0: Foundation (skeleton app)

**Goal:** Dual-pane shell that renders and accepts input. No pi integration yet.

- [ ] Go module initialized (`go mod init github.com/.../pi-mc`)
- [ ] bubbletea app with `tea.WithAltScreen()` (own the terminal)
- [ ] Dual-pane layout: left (35%) + right (65%) with vertical separator
- [ ] Footer F-key bar rendered (non-functional)
- [ ] Tab cycles focus between left/right panes
- [ ] Escape returns focus to right pane
- [ ] Hardcoded mock content in both panes
- [ ] Clean exit on F10 / Ctrl+C

**Deliverable:** Running TUI that shows the MC-style layout with static content.

---

## Phase 1: Pi RPC bridge

**Goal:** Spawn `pi --mode rpc`, parse JSON-Lines, route events.

- [ ] `PiRpcClient` Go struct — spawns `pi --mode rpc --no-session`
- [ ] JSON-Lines parser using `bufio.Scanner` (LF-delimited, ignore `\r`)
- [ ] Typed event structs (`MessageUpdate`, `ToolExecStart`, `QueueUpdate`, etc.)
- [ ] Typed command structs (`Prompt`, `SetModel`, `Abort`, etc.)
- [ ] Command-response correlation via `id` field
- [ ] Graceful subprocess shutdown
- [ ] RPC client test: send `prompt`, receive `message_update` events, print text to stdout

**Deliverable:** Go program that drives pi headlessly and streams text to the terminal.

---

## Phase 2: Chat pane (right)

**Goal:** Replace hardcoded right pane with live pi conversation.

- [ ] Subscribe to `message_update` → append text to chat buffer
- [ ] Handle all delta types: `text_start`, `text_delta`, `text_end`, `thinking_delta`, `toolcall_start`
- [ ] Render user messages (`▸`), assistant messages (`●`), tool results
- [ ] Scrollback engine: ring buffer, viewport, auto-scroll (pause on manual scroll)
- [ ] PgUp/PgDn/Home/End keybindings
- [ ] Text input line at bottom of right pane
- [ ] Enter sends `prompt` RPC command
- [ ] Streaming indicator (cursor blinker at end of streaming message)
- [ ] Shift+Enter for multi-line input
- [ ] Emacs-style line editing (Ctrl+A/E/K/U/W)

**Deliverable:** Full chat experience — type prompt, see streaming response, scroll history.

---

## Phase 3: Dashboard panes (left)

**Goal:** Live Queue, Steps, and Sessions sections.

### 3a: Queue pane
- [ ] Subscribe to `queue_update` events
- [ ] Render steering + follow-up lists
- [ ] Live counter in title bar
- [ ] Distinct icons (`↻` steer, `⏳` follow-up)

### 3b: Steps pane
- [ ] Parse `declare_steps` from tool execution
- [ ] Track `tool_execution_start`/`end` for step status
- [ ] Status icons: `[✓]` `[·]` `[✗]` `[ ]`
- [ ] Color coding (green/red/cyan)
- [ ] Tool count + duration per step

### 3c: Sessions pane
- [ ] Import `SessionManager` from pi SDK for session listing
- [ ] Navigate with `↑↓`, switch on `Enter`
- [ ] Active session marker
- [ ] New session / fork via F-key

**Deliverable:** Full left dashboard — all three sections live and interactive.

---

## Phase 4: Polish

**Goal:** Production-quality experience.

- [ ] Differential rendering for flicker-free streaming
- [ ] ANSI-aware word wrapping (respect escape codes in line width)
- [ ] Visual focus indicators (highlighted pane title, cursor in chat input)
- [ ] Resize handling — recalculate pane dimensions on terminal resize
- [ ] Model cycling (F2) with visual feedback
- [ ] Thinking level cycling (Shift+Tab) with visual feedback
- [ ] F8 manual compaction
- [ ] F10 clean exit (abort agent + close subprocess)
- [ ] Error handling: RPC disconnection, pi crash, parse errors
- [ ] Config file (`~/.pi/agent/pi-mc.json`) for preferences
- [ ] Color theme support (light/dark)

**Deliverable:** Usable daily driver — replace `pi` command with `pi-mc`.

---

## Phase 5: Advanced features

**Goal:** Features that differentiate PI-mc from pi's built-in TUI.

- [ ] File tree pane (F5) — navigable project explorer
- [ ] Conversation tree (F6) — visual session tree with branching
- [ ] External editor integration (Ctrl+G → `$VISUAL`)
- [ ] Vim-style modal editing in chat input
- [ ] `/` command mode — type `/compact`, `/model`, etc.
- [ ] Search in chat history (`/` then search term)
- [ ] Syntax highlighting in code blocks (via chroma or tree-sitter)
- [ ] Image display via Kitty/iTerm2 protocol
- [ ] Session export to HTML/Markdown
- [ ] Multi-agent support (spawn multiple pi subprocesses)

**Deliverable:** Feature-complete agent TUI.

---

## Phase 6: Distribution

**Goal:** Get it in users' hands.

- [ ] `go build` for macOS (ARM + x86), Linux (ARM + x86)
- [ ] Homebrew formula (`brew install pi-mc`)
- [ ] `go install github.com/.../pi-mc@latest`
- [ ] AUR package for Arch Linux
- [ ] Integration test suite
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Release notes + changelog
- [ ] Screenshots + GIF demo for README

---

## Timeline (aspirational)

| Phase | Effort | Cumulative |
|---|---|---|
| 0: Foundation | 1-2 days | 1-2 days |
| 1: RPC bridge | 2-3 days | 3-5 days |
| 2: Chat pane | 3-5 days | 6-10 days |
| 3: Dashboard panes | 3-5 days | 9-15 days |
| 4: Polish | 5-7 days | 14-22 days |
| 5: Advanced | 7-14 days | 21-36 days |
| 6: Distribution | 3-5 days | 24-41 days |

**Total: ~4-6 weeks** of focused work for a polished v1.0.

---

## Immediate next step

**Phase 0, Step 1:** Initialize Go module, create `main.go` with bubbletea dual-pane layout, hardcoded content, working F-key bar and Tab focus. This validates the visual concept before touching any pi integration.
