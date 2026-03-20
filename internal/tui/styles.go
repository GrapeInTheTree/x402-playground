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

// LayoutPage places body content at the top and hints anchored at the bottom,
// filling the full height. Content is horizontally centered within the width.
func LayoutPage(body, hints string, width, height int) string {
	centeredBody := lipgloss.PlaceHorizontal(width, lipgloss.Center, body)

	bodyH := lipgloss.Height(centeredBody)
	hintsH := lipgloss.Height(hints)
	gap := max(height-bodyH-hintsH, 0)

	return centeredBody + "\n" +
		lipgloss.NewStyle().Height(gap).Render("") +
		MutedStyle.Render(hints)
}
