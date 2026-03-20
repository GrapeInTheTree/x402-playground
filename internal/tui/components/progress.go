package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// Progress renders a step progress indicator like ●●●○○○○○○○
type Progress struct {
	Total   int
	Current int
}

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
	label := lipgloss.NewStyle().
		Foreground(tui.ColorAccent).
		Bold(true).
		Render(strings.Replace("Step X/Y", "X", string(rune('0'+current%10)), 1))

	// Use Sprintf for proper formatting
	stepText := lipgloss.NewStyle().
		Foreground(tui.ColorAccent).
		Bold(true).
		Render(stepString(current, total))

	desc := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#D1D5DB")).
		Render(description)

	_ = label
	return stepText + " " + desc
}

func stepString(current, total int) string {
	digits := "0123456789"
	if current > 9 || total > 9 {
		return string(rune('0'+current/10)) + string(digits[current%10]) + "/" +
			string(rune('0'+total/10)) + string(digits[total%10])
	}
	return string(digits[current]) + "/" + string(digits[total])
}
