package practice

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-demo/internal/demo"
	"github.com/GrapeInTheTree/x402-demo/internal/tui"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/components"
)

// StepInfo holds the display data for a single step in a panel.
type StepInfo struct {
	Actor   string // "Client", "Resource", "Facilitator"
	Action  string // brief description
	Detail  string // detailed content (JSON, etc.)
	Status  string // "pending", "running", "done", "error"
}

// RenderFlowPanels renders the 3-column practice view.
func RenderFlowPanels(
	step int, totalSteps int,
	clientSteps, resourceSteps, facilitatorSteps []StepInfo,
	width int,
) string {
	progress := components.Progress{Total: totalSteps, Current: step}
	stepDesc := demo.StepDescription(step + 1)
	stepLabel := fmt.Sprintf("Step %d/%d: %s", step+1, totalSteps, stepDesc)

	header := lipgloss.JoinHorizontal(lipgloss.Center,
		"    ",
		progress.View(),
		"  ",
		lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).Render(stepLabel),
	)

	colWidth := (width - 8) / 3
	if colWidth < 25 {
		colWidth = 25
	}

	clientPanel := renderStepPanel("Client", clientSteps, step, colWidth, tui.ColorSecondary)
	resourcePanel := renderStepPanel("Resource Server", resourceSteps, step, colWidth, tui.ColorPrimary)
	facilitatorPanel := renderStepPanel("Facilitator", facilitatorSteps, step, colWidth, tui.ColorAccent)

	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		clientPanel,
		" ",
		resourcePanel,
		" ",
		facilitatorPanel,
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		panels,
	)
}

func renderStepPanel(title string, steps []StepInfo, currentStep int, width int, color lipgloss.Color) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Width(width - 2).
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)

	var b strings.Builder
	b.WriteString(titleStyle.Render(title) + "\n")

	for i, s := range steps {
		if i > currentStep+1 {
			break // Don't show future steps
		}

		var icon string
		var lineStyle lipgloss.Style
		switch s.Status {
		case "done":
			icon = "✓"
			lineStyle = lipgloss.NewStyle().Foreground(tui.ColorSuccess)
		case "running":
			icon = "►"
			lineStyle = lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true)
		case "error":
			icon = "✗"
			lineStyle = lipgloss.NewStyle().Foreground(tui.ColorError)
		default:
			icon = "○"
			lineStyle = lipgloss.NewStyle().Foreground(tui.ColorMuted)
		}

		b.WriteString(lineStyle.Render(fmt.Sprintf(" %s %s", icon, s.Action)) + "\n")

		if s.Status == "running" && s.Detail != "" {
			detail := lipgloss.NewStyle().
				Foreground(tui.ColorMuted).
				Width(width - 8).
				Render("   " + s.Detail)
			b.WriteString(detail + "\n")
		}
	}

	return style.Render(b.String())
}
