package ui

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileTreeModel struct {
	width  int
	height int
	active bool

	currentDir string
	files      []os.DirEntry
	cursor     int
}

func NewFileTree() FileTreeModel {
	dir, _ := os.Getwd()
	m := FileTreeModel{currentDir: dir}
	m.loadDir()
	return m
}

func (m *FileTreeModel) loadDir() {
	entries, _ := os.ReadDir(m.currentDir)
	m.files = entries
	m.cursor = 0
}

func (m FileTreeModel) Init() tea.Cmd {
	return nil
}

func (m FileTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.active {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.files)-1 {
					m.cursor++
				}
			case "enter":
				// If a directory is selected, we could navigate. For now just placeholder.
			}
		}
	}
	return m, nil
}

func (m FileTreeModel) View() string {
	style := paneStyle
	if m.active {
		style = activePaneStyle
	}

	innerH := m.height - 2
	innerW := m.width - 4
	if innerW < 0 { innerW = 0 }
	if innerH < 0 { innerH = 0 }

	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Render("📁 " + m.currentDir + "\n"))
	
	start := 0
	if m.cursor >= innerH-1 {
		start = m.cursor - (innerH - 2)
	}
	
	for i := start; i < len(m.files) && i < start+innerH-1; i++ {
		f := m.files[i]
		prefix := "  "
		if i == m.cursor {
			prefix = "▸ "
		}
		name := f.Name()
		if f.IsDir() {
			name += "/"
		}
		sb.WriteString(prefix)
		sb.WriteString(name)
		sb.WriteString("\n")
	}

	lines := strings.Split(sb.String(), "\n")
	if len(lines) > innerH {
		lines = lines[:innerH]
	}
	content := strings.Join(lines, "\n")
	
	// Apply inner constraints and explicit background
	content = lipgloss.NewStyle().
		Background(CurrentTheme.Background).
		MaxWidth(innerW).
		Render(content)

	return style.Width(innerW).Height(innerH).Render(content)
}

func (m *FileTreeModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *FileTreeModel) SetActive(active bool) {
	m.active = active
}
