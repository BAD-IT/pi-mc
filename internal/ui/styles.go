package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Border color for the inactive pane
	inactiveBorderColor = lipgloss.Color("238")
	// Border color for the active pane
	activeBorderColor = lipgloss.Color("62") // Indigo

	paneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(inactiveBorderColor).
		Padding(1, 2)

	activePaneStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(activeBorderColor).
		Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
		MarginTop(0).
		Padding(0, 1)
	
	fKeyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("238")).
		Padding(0, 1).
		MarginRight(1)
	
	fKeyLabelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))
)
