package ui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Background       lipgloss.AdaptiveColor
	TextColor        lipgloss.AdaptiveColor
	PaneBorder       lipgloss.AdaptiveColor
	ActivePaneBorder lipgloss.AdaptiveColor
	FooterBg         lipgloss.AdaptiveColor
	FooterFg         lipgloss.AdaptiveColor
	FKeyBg           lipgloss.AdaptiveColor
	FKeyFg           lipgloss.AdaptiveColor
}

// True Midnight Commander classic DOS theme
var mcTheme = Theme{
	Background:       lipgloss.AdaptiveColor{Light: "#0000AA", Dark: "#0000AA"}, // Classic DOS Blue
	TextColor:        lipgloss.AdaptiveColor{Light: "#AAAAAA", Dark: "#AAAAAA"}, // Light Gray
	PaneBorder:       lipgloss.AdaptiveColor{Light: "#55FFFF", Dark: "#55FFFF"}, // Cyan
	ActivePaneBorder: lipgloss.AdaptiveColor{Light: "#FFFF55", Dark: "#FFFF55"}, // Bright Yellow
	FooterBg:         lipgloss.AdaptiveColor{Light: "0", Dark: "0"},             // Black
	FooterFg:         lipgloss.AdaptiveColor{Light: "255", Dark: "255"},         // White
	FKeyBg:           lipgloss.AdaptiveColor{Light: "0", Dark: "0"},             // Black
	FKeyFg:           lipgloss.AdaptiveColor{Light: "#55FFFF", Dark: "#55FFFF"}, // Cyan
}

// Light theme as an alternative
var lightTheme = Theme{
	Background:       lipgloss.AdaptiveColor{Light: "255", Dark: "255"},
	TextColor:        lipgloss.AdaptiveColor{Light: "235", Dark: "235"},
	PaneBorder:       lipgloss.AdaptiveColor{Light: "250", Dark: "250"},
	ActivePaneBorder: lipgloss.AdaptiveColor{Light: "27", Dark: "27"},
	FooterBg:         lipgloss.AdaptiveColor{Light: "235", Dark: "235"},
	FooterFg:         lipgloss.AdaptiveColor{Light: "255", Dark: "255"},
	FKeyBg:           lipgloss.AdaptiveColor{Light: "235", Dark: "235"},
	FKeyFg:           lipgloss.AdaptiveColor{Light: "255", Dark: "255"},
}

var CurrentTheme = mcTheme

var (
	basePaneStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			Padding(0, 1)

	baseFooterStyle = lipgloss.NewStyle().
			MarginTop(0).
			Padding(0, 1)

	baseFKeyStyle = lipgloss.NewStyle().
			Padding(0, 1)

	baseFKeyLabelStyle = lipgloss.NewStyle()
)

// Computed styles used across components
var (
	paneStyle       lipgloss.Style
	activePaneStyle lipgloss.Style
	footerStyle     lipgloss.Style
	fKeyStyle       lipgloss.Style
	fKeyLabelStyle  lipgloss.Style
	appBgStyle      lipgloss.Style
)

func applyTheme(theme string) {
	if theme == "light" {
		CurrentTheme = lightTheme
	} else {
		CurrentTheme = mcTheme
	}

	appBgStyle = lipgloss.NewStyle().
		Background(CurrentTheme.Background).
		Foreground(CurrentTheme.TextColor)

	paneStyle = basePaneStyle.
		BorderForeground(CurrentTheme.PaneBorder).
		Background(CurrentTheme.Background).
		Foreground(CurrentTheme.TextColor)

	activePaneStyle = basePaneStyle.
		BorderForeground(CurrentTheme.ActivePaneBorder).
		Background(CurrentTheme.Background).
		Foreground(CurrentTheme.TextColor)

	footerStyle = baseFooterStyle.
		Background(CurrentTheme.FooterBg).
		Foreground(CurrentTheme.FooterFg)

	fKeyStyle = baseFKeyStyle.
		Background(CurrentTheme.FKeyBg).
		Foreground(CurrentTheme.FKeyFg)

	fKeyLabelStyle = baseFKeyLabelStyle.
		Background(CurrentTheme.FooterBg).
		Foreground(CurrentTheme.FooterFg)
}
