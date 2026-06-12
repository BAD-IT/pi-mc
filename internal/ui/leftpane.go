package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LeftPaneModel struct {
	width  int
	height int
	active bool
}

func NewLeftPane() LeftPaneModel {
	return LeftPaneModel{}
}

func (m LeftPaneModel) Init() tea.Cmd {
	return nil
}

func (m LeftPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m LeftPaneModel) View() string {
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

	content := "📋 Queue\n\n" +
		"· Audit error paths\n" +
		"· Add tests\n\n\n" +
		"📊 Steps 2/3\n\n" +
		"1 [✓] Read auth module\n" +
		"2 [·] Propose refactor\n" +
		"3 [ ] Apply changes\n\n\n" +
		"🗂 Sessions (4)\n\n" +
		"▸ refactor-auth (active)\n" +
		"  fix-login-bug\n" +
		"  rate-limiting\n"

	return style.Width(innerW).Height(innerH).Render(content)
}

func (m *LeftPaneModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *LeftPaneModel) SetActive(active bool) {
	m.active = active
}
