# PI-mc

**Midnight Commander-style full-screen TUI for the pi coding agent.**

PI-mc replaces pi's line-oriented chat with a dual-pane dashboard: orchestration controls on the left, streaming conversation on the right, F-key action bar at the bottom. Everything keyboard-driven, always visible.

```
┌─ Dashboard ─────────────────┬─ Chat ───────────────────────────────┐
│                             │                                      │
│  📋 Queue                   │  smith: refactor the auth module     │
│  ┌──────────────────────┐   │                                      │
│  │ · Audit error paths  │   │  pi: Let me read the current auth   │
│  │ · Add tests          │   │      module and propose changes.     │
│  └──────────────────────┘   │                                      │
│                             │    [✓] Read auth module    3t, 2s    │
│  📊 Steps 2/3               │    [·] Propose refactor     1t      │
│  1 [✓] Read auth module     │    [ ] Apply changes          —      │
│  2 [·] Propose refactor     │                                      │
│  3 [ ] Apply changes        │  ─────────────────────────────────── │
│                             │  ▸ Type a message, Enter to send     │
│  🗂 Sessions (4)            │                                      │
│  ▸ refactor-auth (active)   │                                      │
│    fix-login-bug            │                                      │
│    rate-limiting            │                                      │
│                             │                                      │
├─────────────────────────────┴──────────────────────────────────────┤
│ 1Queue  2Model  3Steps  4Sessions  5Files  6Tree  8Compact  10Quit │
└────────────────────────────────────────────────────────────────────┘
```

## How it works

PI-mc spawns `pi --mode rpc` as a headless subprocess and drives it via the JSON-Lines RPC protocol. All AI logic, tool execution, and session management stays in pi. PI-mc is pure rendering + input.

```
┌──────────┐  stdin: commands  ┌──────────┐
│  PI-mc  │ ◄──────────────── │ pi --rpc │
│  (Go)    │ ────────────────► │ (headless)│
└──────────┘  stdout: events   └──────────┘
```

## Project status

**Pre-alpha — documentation and planning phase.** No code yet.

See [roadmap.md](docs/roadmap.md) for phased milestones.

## Design

- [Concept](docs/concept.md) — vision, philosophy, architecture
- [Features](docs/features.md) — complete feature catalogue
- [Roadmap](docs/roadmap.md) — phased development plan

## Tech stack

| Layer | Choice |
|---|---|
| Language | Go |
| TUI framework | [bubbletea](https://github.com/charmbracelet/bubbletea) |
| Styling | [lipgloss](https://github.com/charmbracelet/lipgloss) |
| Agent backend | `pi --mode rpc` (JSON-Lines protocol) |
| Distribution | Single static binary (macOS + Linux) |

## Inspiration

- [Midnight Commander](https://midnight-commander.org/) — dual-pane layout, F-key bar
- [lazygit](https://github.com/jesseduffield/lazygit) — bubbletea TUI for git
- [k9s](https://github.com/derailed/k9s) — terminal dashboard for Kubernetes

## License

MIT
