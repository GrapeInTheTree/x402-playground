package components

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// Panel renders a bordered box with a title and content.
type Panel struct {
	Title   string
	Content string
	Status  string
	Width   int
	Active  bool
}

func (p Panel) View() string {
	style := tui.BorderStyle.Width(p.Width - 4)
	if p.Active {
		style = tui.ActiveBorderStyle.Width(p.Width - 4)
	}

	titleStyle := tui.SubtitleStyle
	if p.Active {
		titleStyle = lipgloss.NewStyle().
			Foreground(tui.ColorPrimary).
			Bold(true)
	}

	content := titleStyle.Render(p.Title) + "\n\n" + p.Content

	if p.Status != "" {
		content += "\n\n" + tui.MutedStyle.Render(p.Status)
	}

	return style.Render(content)
}
