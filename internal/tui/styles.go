package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#7C3AED") // violet
	ColorSecondary = lipgloss.Color("#06B6D4") // cyan
	ColorAccent    = lipgloss.Color("#F59E0B") // amber
	ColorSuccess   = lipgloss.Color("#10B981") // green
	ColorError     = lipgloss.Color("#EF4444") // red
	ColorMuted     = lipgloss.Color("#6B7280") // gray
	ColorBorder    = lipgloss.Color("#374151") // dark gray
	ColorBg        = lipgloss.Color("#111827") // near-black
	ColorSubtle    = lipgloss.Color("#1F2937") // slightly lighter than bg
	ColorHighlight = lipgloss.Color("#2D1B69") // subtle violet bg for selection

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	AccentStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	// Layout
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	ActiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(1, 2)

	// Menu
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(ColorPrimary).
				Bold(true).
				SetString("> ")

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)
)

// LayoutPage renders a full-screen layout with a colored header bar at top,
// content filling the middle, and a status bar pinned at the bottom.
func LayoutPage(body, hints string, width, height int) string {
	if width <= 0 || height <= 0 {
		return body
	}

	// Status bar: full-width background with colored segments
	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Background(ColorSubtle).
		Width(width).
		Padding(0, 2).
		Render("  " + hints)

	statusH := lipgloss.Height(statusBar)

	// Content area: fill remaining height with left-aligned padded content
	contentH := max(height-statusH, 1)
	content := lipgloss.NewStyle().
		Width(width).
		Height(contentH).
		Padding(0, 3).
		Render(body)

	return lipgloss.JoinVertical(lipgloss.Top, content, statusBar)
}

// LayoutPageCentered renders body inside a rounded border card, centered both
// vertically and horizontally, with a status bar pinned at the bottom.
func LayoutPageCentered(body, hints string, width, height int) string {
	if width <= 0 || height <= 0 {
		return body
	}

	// Wrap body in a bordered card
	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 3).
		Render(body)

	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Background(ColorSubtle).
		Width(width).
		Padding(0, 2).
		Render("  " + hints)

	statusH := lipgloss.Height(statusBar)
	bodyAreaH := max(height-statusH, 1)

	centered := lipgloss.Place(width, bodyAreaH,
		lipgloss.Center, lipgloss.Center,
		card)

	return lipgloss.JoinVertical(lipgloss.Top, centered, statusBar)
}
