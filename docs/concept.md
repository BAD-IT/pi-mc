# PI-mc — Concept

## What is PI-mc?

PI-mc is a **full-screen terminal user interface** for the [pi coding agent](https://github.com/earendil-works/pi-mono). It replaces pi's line-oriented chat interface with a Midnight Commander-style dual-pane dashboard where every feature is visible, keyboard-accessible, and always in reach.

## Why

Pi today is a **REPL**: prompt at the bottom, history scrolls up. Powerful, but everything is hidden behind `/commands`, keybindings you have to memorize, and context you can't see at a glance.

PI-mc makes the agent's state **spatially persistent**:

| Pi today (line-oriented) | PI-mc (screen-first) |
|---|---|
| Progress steps flash briefly | Steps live in left pane, always visible |
| Queue is invisible | Pending steering/follow-up shown in list |
| Sessions hidden behind `/tree` | Session list navigable with arrow keys |
| One input line, one prompt | Chat + multi-pane dashboard |
| Everything is sequential scrollback | Left pane = control, right = output |

## Philosophy

### 1. Midnight Commander as north star

Two panes. F-key bar at the bottom. Everything keyboard-driven. Tab switches focus. MC proves this layout works for complex data — PI-mc applies it to agent conversations.

### 2. The left pane is the cockpit

Not a file tree. An **orchestration dashboard** showing:
- What's queued (steering/follow-up messages)
- What's in progress (work steps with live status)
- What sessions exist (fork/switch with a keystroke)

The left pane answers "what can I do?" — the right pane answers "what's happening?"

### 3. Zero memorization

F-keys are labeled. Every pane has a visible title. Selection is highlighted. If you can see it, you can navigate to it.

### 4. Keyboard purity

No mouse. No drag-resize. No animations. Just `Tab` to switch focus, `↑↓` to navigate, `Enter` to act. Same as MC.

### 5. Pi as headless engine

PI-mc doesn't re-implement pi. It spawns `pi --mode rpc` as a subprocess and drives it via the JSON-Lines protocol. All AI logic, tool execution, and session management stays in pi. PI-mc is **pure rendering + input**.

## Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                        PI-mc (Go + bubbletea)                   │
│                                                                  │
│  ┌─────────────────────┐  ┌──────────────────────────────────┐  │
│  │    Left Pane (35%)  │  │      Right Pane (65%)            │  │
│  │                     │  │                                  │  │
│  │  ┌─ Queue ────────┐ │  │  ┌────────────────────────────┐ │  │
│  │  │ pending msgs   │ │  │  │                            │ │  │
│  │  │ steer / follow │ │  │  │  Chat + streaming messages │ │  │
│  │  └────────────────┘ │  │  │                            │ │  │
│  │                     │  │  │  smith: refactor auth      │ │  │
│  │  ┌─ Steps ────────┐ │  │  │  pi:   I'll read the      │ │  │
│  │  │ 1 [✓] Read     │ │  │  │         module first...    │ │  │
│  │  │ 2 [·] Propose  │ │  │  │                            │ │  │
│  │  │ 3 [ ] Apply    │ │  │  │  [✓] read auth.ts   3t,2s │ │  │
│  │  └────────────────┘ │  │  │  [·] propose        1t     │ │  │
│  │                     │  │  │  [ ] apply            —     │ │  │
│  │  ┌─ Sessions ─────┐ │  │  └────────────────────────────┘ │  │
│  │  │ ▸ refactor     │ │  │                                  │  │
│  │  │   fix-login    │ │  │  ▸ Type a message, Enter to   │  │
│  │  │   rate-limit   │ │  │    send                        │  │
│  │  └────────────────┘ │  │                                  │  │
│  └─────────────────────┘  └──────────────────────────────────┘  │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │ 1Queue  2Model  3Steps  4Sessions  5Files  6Tree  10Quit  │  │
│  └────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
                           │
                           │ stdin: JSON commands
                           │ stdout: JSON events (text_delta, tool_exec, queue_update)
                           ▼
                    ┌──────────────┐
                    │  pi --mode   │
                    │  rpc         │
                    │  (subprocess)│
                    └──────────────┘
```

## RPC protocol

PI-mc communicates with pi exclusively via the [JSON-Lines RPC protocol](https://github.com/earendil-works/pi-mono/blob/main/docs/rpc.md):

- **Commands** (PI-mc → pi): `prompt`, `steer`, `follow_up`, `abort`, `set_model`, `set_thinking_level`, `compact`, `bash`, `switch_session`, `fork`, `clone`
- **Events** (pi → PI-mc): `message_update` (text_delta, thinking_delta, toolcall_delta), `tool_execution_start`, `tool_execution_end`, `queue_update`, `agent_start`, `agent_end`, `turn_start`, `turn_end`, `compaction_start`, `compaction_end`
- **Extension UI** (bidirectional): `select`, `confirm`, `input`, `editor`, `notify`, `setStatus`, `setWidget`

Every pane subscribes to the events it needs. The render loop is driven by event arrival + user input.

## Technology choices

| Choice | Rationale |
|---|---|
| **Go** | Single static binary, fast compiles, excellent TUI ecosystem |
| **bubbletea** | Elm Architecture (Model → Update → View), perfect for event-driven UI |
| **lipgloss** | MC-style styling (reverse-video, borders, dimmed text) |
| **bufio.Scanner** | JSON-Lines parsing is ~50 lines in Go |
| **pi --mode rpc** | Zero reimplementation — all agent logic stays in pi |
