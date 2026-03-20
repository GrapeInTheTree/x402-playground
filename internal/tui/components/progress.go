package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// Progress renders a step progress indicator like ●●●○○○○○○○
type Progress struct {
	Total   int
	Current int
}

// View renders the progress indicator as filled and empty circles.
func (p Progress) View() string {
	var b strings.Builder

	filled := lipgloss.NewStyle().Foreground(tui.ColorPrimary)
	empty := lipgloss.NewStyle().Foreground(tui.ColorBorder)

	for i := 0; i < p.Total; i++ {
		if i < p.Current {
			b.WriteString(filled.Render("●"))
		} else if i == p.Current {
			b.WriteString(lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).Render("●"))
		} else {
			b.WriteString(empty.Render("○"))
		}
	}

	return b.String()
}

// StepLabel returns a formatted step label like "Step 3/10: Description"
func StepLabel(current, total int, description string) string {
	stepText := lipgloss.NewStyle().
		Foreground(tui.ColorAccent).
		Bold(true).
		Render(fmt.Sprintf("%d/%d", current, total))

	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D1D5DB")).
		Render(description)

	return stepText + " " + desc
}
