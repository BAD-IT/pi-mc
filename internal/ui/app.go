package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppModel struct {
	leftPane   LeftPaneModel
	rightPane  RightPaneModel
	footer     FooterModel
	
	// focus state: false = left, true = right
	focusRight bool
	
	width  int
	height int
}

func NewApp() AppModel {
	return AppModel{
		leftPane:   NewLeftPane(),
		rightPane:  NewRightPane(),
		footer:     NewFooter(),
		focusRight: true, // Right pane (Chat) is default
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "f10":
			return m, tea.Quit
		case "tab":
			m.focusRight = !m.focusRight
		case "esc":
			m.focusRight = true
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Recalculate panes. Let's allocate 2 lines for footer (content + margin)
		paneHeight := m.height - 2
		if paneHeight < 0 {
			paneHeight = 0
		}
		
		leftWidth := int(float64(m.width) * 0.35)
		rightWidth := m.width - leftWidth
		
		m.leftPane.SetSize(leftWidth, paneHeight)
		m.rightPane.SetSize(rightWidth, paneHeight)
		
		f, _ := m.footer.Update(msg)
		m.footer = f.(FooterModel)
	}
	
	m.leftPane.SetActive(!m.focusRight)
	m.rightPane.SetActive(m.focusRight)
	
	return m, nil
}

func (m AppModel) View() string {
	panes := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.leftPane.View(),
		m.rightPane.View(),
	)
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		panes,
		m.footer.View(),
	)
}
