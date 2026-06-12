package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/earendil-works/pi-mc/internal/rpc"
)

type StepState int

const (
	StepPending StepState = iota
	StepRunning
	StepDone
)

type StepInfo struct {
	ID          int
	Description string
	State       StepState
	ToolCount   int
}

type SessionInfo struct {
	Name   string
	Active bool
}

type LeftPaneModel struct {
	width  int
	height int
	active bool

	queueSteering []string
	queueFollowUp []string

	steps       []StepInfo
	currentStep int

	sessions     []SessionInfo
	sessionIndex int
}

func NewLeftPane() LeftPaneModel {
	return LeftPaneModel{
		sessions: []SessionInfo{
			{Name: "refactor-auth", Active: true},
			{Name: "fix-login-bug"},
			{Name: "rate-limiting"},
		},
	}
}

func (m LeftPaneModel) Init() tea.Cmd {
	return nil
}

func (m LeftPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.active {
			switch msg.String() {
			case "up", "k":
				if m.sessionIndex > 0 {
					m.sessionIndex--
				}
			case "down", "j":
				if m.sessionIndex < len(m.sessions)-1 {
					m.sessionIndex++
				}
			case "enter":
				for i := range m.sessions {
					m.sessions[i].Active = (i == m.sessionIndex)
				}
			}
		}

	case rpcEventMsg:
		switch msg.Type {
		case "queue_update":
			var q rpc.QueueUpdate
			if err := json.Unmarshal(msg.Payload, &q); err == nil {
				m.queueSteering = q.Steering
				m.queueFollowUp = q.FollowUp
			}
		case "declare_steps":
			var d rpc.DeclareSteps
			if err := json.Unmarshal(msg.Payload, &d); err == nil {
				m.steps = make([]StepInfo, len(d.Steps))
				for i, s := range d.Steps {
					m.steps[i] = StepInfo{
						ID:          s.ID,
						Description: s.Description,
						State:       StepPending,
					}
				}
				m.currentStep = 0
				if len(m.steps) > 0 {
					m.steps[0].State = StepRunning
				}
			}
		case "tool_execution_start":
			if m.currentStep < len(m.steps) {
				m.steps[m.currentStep].ToolCount++
			}
		case "tool_execution_end":
			if m.currentStep < len(m.steps) {
				m.steps[m.currentStep].State = StepDone
				m.currentStep++
				if m.currentStep < len(m.steps) {
					m.steps[m.currentStep].State = StepRunning
				}
			}
		}
	}

	return m, nil
}

func (m LeftPaneModel) View() string {
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

	var sb strings.Builder

	// Queue Section
	qCount := len(m.queueSteering) + len(m.queueFollowUp)
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("📋 Queue (%d)\n", qCount)))
	for _, s := range m.queueSteering {
		sb.WriteString(fmt.Sprintf("↻ %s\n", s))
	}
	for _, f := range m.queueFollowUp {
		sb.WriteString(fmt.Sprintf("⏳ %s\n", f))
	}
	if qCount == 0 {
		sb.WriteString("  (empty)\n")
	}
	sb.WriteString("\n")

	// Steps Section
	completed := 0
	for _, s := range m.steps {
		if s.State == StepDone {
			completed++
		}
	}
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("📊 Steps %d/%d\n", completed, len(m.steps))))
	for _, s := range m.steps {
		icon := "[ ]"
		switch s.State {
		case StepRunning:
			icon = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render("[·]")
		case StepDone:
			icon = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("[✓]")
		}
		
		toolStr := ""
		if s.ToolCount > 0 {
			toolStr = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(fmt.Sprintf(" %dt", s.ToolCount))
		}
		sb.WriteString(fmt.Sprintf("%d %s %s%s\n", s.ID, icon, s.Description, toolStr))
	}
	if len(m.steps) == 0 {
		sb.WriteString("  (no active task)\n")
	}
	sb.WriteString("\n")

	// Sessions Section
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("🗂 Sessions (%d)\n", len(m.sessions))))
	for i, s := range m.sessions {
		prefix := "  "
		if i == m.sessionIndex {
			prefix = "▸ "
		}
		
		name := s.Name
		if s.Active {
			name = lipgloss.NewStyle().Bold(true).Render(name + " (active)")
		} else if i == m.sessionIndex {
			name = lipgloss.NewStyle().Underline(true).Render(name)
		}
		
		sb.WriteString(fmt.Sprintf("%s%s\n", prefix, name))
	}

	return style.Width(innerW).Height(innerH).Render(sb.String())
}

func (m *LeftPaneModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *LeftPaneModel) SetActive(active bool) {
	m.active = active
}
