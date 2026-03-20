package components

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-demo/internal/tui"
)

// StatusBar renders a bottom status bar with navigation hints.
type StatusBar struct {
	Width int
}

func (s StatusBar) View(hints string) string {
	style := lipgloss.NewStyle().
		Foreground(tui.ColorMuted).
		Width(s.Width)

	return style.Render(hints)
}
