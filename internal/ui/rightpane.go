package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PromptMsg string

type RightPaneModel struct {
	viewport textarea.Model // Actually we use viewport.Model, see below
	
	vp       viewport.Model
	ta       textarea.Model
	
	width  int
	height int
	active bool
	
	messages             []string
	isAssistantStreaming bool
}

func NewRightPane() RightPaneModel {
	ta := textarea.New()
	ta.Placeholder = "Type a message, Enter to send..."
	ta.Focus()
	ta.Prompt = "▸ "
	ta.CharLimit = 1000
	ta.SetWidth(30)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	vp := viewport.New(30, 10)

	return RightPaneModel{
		ta:       ta,
		vp:       vp,
		messages: []string{},
	}
}

func (m RightPaneModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m RightPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.active && msg.Type == tea.KeyEnter {
			v := m.ta.Value()
			if strings.TrimSpace(v) != "" {
				m.ta.Reset()
				return m, func() tea.Msg {
					return PromptMsg(v)
				}
			}
			return m, nil
		}
	}

	if m.active {
		m.ta, tiCmd = m.ta.Update(msg)
		cmds = append(cmds, tiCmd)
		
		m.vp, vpCmd = m.vp.Update(msg)
		cmds = append(cmds, vpCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m RightPaneModel) View() string {
	style := paneStyle
	if m.active {
		style = activePaneStyle
	}

	innerH := m.height - 2
	if innerH < 0 {
		innerH = 0
	}
	innerW := m.width - 4
	if innerW < 0 {
		innerW = 0
	}

	vpView := m.vp.View()
	taView := m.ta.View()

	sep := strings.Repeat("─", innerW)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		vpView,
		sep,
		taView,
	)

	return style.Width(innerW).Height(innerH).Render(content)
}

func (m *RightPaneModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	innerH := height - 2
	innerW := width - 4
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}

	m.ta.SetWidth(innerW)
	
	taHeight := 3
	vpHeight := innerH - taHeight - 1
	if vpHeight < 0 {
		vpHeight = 0
	}

	m.vp.Width = innerW
	m.vp.Height = vpHeight
}

func (m *RightPaneModel) SetActive(active bool) {
	m.active = active
	if active {
		m.ta.Focus()
	} else {
		m.ta.Blur()
	}
}

func (m *RightPaneModel) AddMessage(msg string) {
	m.messages = append(m.messages, msg)
	m.isAssistantStreaming = false
	m.vp.SetContent(strings.Join(m.messages, "\n\n"))
	m.vp.GotoBottom()
}

func (m *RightPaneModel) AppendStream(chunk string) {
	if !m.isAssistantStreaming {
		m.messages = append(m.messages, "● pi: "+chunk)
		m.isAssistantStreaming = true
	} else {
		m.messages[len(m.messages)-1] += chunk
	}
	m.vp.SetContent(strings.Join(m.messages, "\n\n"))
	m.vp.GotoBottom()
}
