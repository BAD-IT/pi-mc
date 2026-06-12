package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RightPaneModel struct {
	width  int
	height int
	active bool
}

func NewRightPane() RightPaneModel {
	return RightPaneModel{}
}

func (m RightPaneModel) Init() tea.Cmd {
	return nil
}

func (m RightPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m RightPaneModel) View() string {
	style := paneStyle
	if m.active {
		style = activePaneStyle
	}

	// Calculate inner dimensions
	innerH := m.height - 2
	if innerH < 0 {
		innerH = 0
	}
	innerW := m.width - 4
	if innerW < 0 {
		innerW = 0
	}

	content := "▸ smith: refactor the auth module\n\n" +
		"● pi: Let me read the current auth\n" +
		"      module and propose changes.\n\n" +
		"  [✓] Read auth module    3t, 2s\n" +
		"  [·] Propose refactor     1t\n" +
		"  [ ] Apply changes          —\n\n" +
		"───────────────────────────────────\n" +
		"▸ Type a message, Enter to send\n"

	return style.Width(innerW).Height(innerH).Render(content)
}

func (m *RightPaneModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *RightPaneModel) SetActive(active bool) {
	m.active = active
}
