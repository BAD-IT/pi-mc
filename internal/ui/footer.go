package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type FooterModel struct {
	width int
}

func NewFooter() FooterModel {
	return FooterModel{}
}

func (m FooterModel) Init() tea.Cmd {
	return nil
}

func (m FooterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

func (m FooterModel) View() string {
	keys := []struct {
		num   int
		label string
	}{
		{1, "Help"},
		{2, "Model"},
		{3, "Steps"},
		{4, "Sessions"},
		{5, "Tree"},
		{6, "Convo"},
		{7, "Search"},
		{8, "Compact"},
		{9, "Menu"},
		{10, "Quit"},
	}

	var sb strings.Builder
	for _, k := range keys {
		fkey := fKeyStyle.Render(fmt.Sprintf("F%d", k.num))
		label := fKeyLabelStyle.Render(k.label)
		sb.WriteString(fkey)
		sb.WriteString(label)
		sb.WriteString(" ")
	}

	return footerStyle.Width(m.width).Render(sb.String())
}
