package practice

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/demo"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

const (
	inactiveAction = "—"
	minPanelWidth  = 25
	panelGap       = 8
)

// StepInfo holds the display data for a single step in a panel.
type StepInfo struct {
	Actor  string // "Client", "Resource", "Facilitator"
	Action string
	Detail string
	Status string // "pending", "running", "done", "error"
}

// stepDesc holds the 3-actor description for initializing steps.
type stepDesc struct {
	Client, Resource, Facilitator string
}

// stepState groups step info for all three actors.
type stepState struct {
	client      StepInfo
	resource    StepInfo
	facilitator StepInfo
}

// stepManager encapsulates the shared step navigation logic.
type stepManager struct {
	flow  *demo.FlowState
	steps [10]stepState
}

func newStepManager(flow *demo.FlowState, descriptions []stepDesc) stepManager {
	sm := stepManager{flow: flow}
	for i, d := range descriptions {
		sm.steps[i] = stepState{
			client:      StepInfo{Actor: "Client", Action: d.Client, Status: "pending"},
			resource:    StepInfo{Actor: "Resource", Action: d.Resource, Status: "pending"},
			facilitator: StepInfo{Actor: "Facilitator", Action: d.Facilitator, Status: "pending"},
		}
	}
	if len(descriptions) > 0 {
		sm.steps[0].client.Status = "running"
	}
	return sm
}

func (sm *stepManager) next() {
	if sm.flow.CurrentStep >= sm.flow.TotalSteps-1 {
		return
	}
	s := &sm.steps[sm.flow.CurrentStep]
	s.client.Status = "done"
	s.resource.Status = "done"
	s.facilitator.Status = "done"

	sm.flow.CurrentStep++
	sm.markRunning(&sm.steps[sm.flow.CurrentStep])
}

func (sm *stepManager) prev() {
	if sm.flow.CurrentStep <= 0 {
		return
	}
	s := &sm.steps[sm.flow.CurrentStep]
	s.client.Status = "pending"
	s.resource.Status = "pending"
	s.facilitator.Status = "pending"

	sm.flow.CurrentStep--
	sm.markRunning(&sm.steps[sm.flow.CurrentStep])
}

func (sm *stepManager) markRunning(s *stepState) {
	if s.client.Action != inactiveAction {
		s.client.Status = "running"
	}
	if s.resource.Action != inactiveAction {
		s.resource.Status = "running"
	}
	if s.facilitator.Action != inactiveAction {
		s.facilitator.Status = "running"
	}
}

func (sm *stepManager) markStepDone() {
	if sm.flow.CurrentStep >= sm.flow.TotalSteps {
		return
	}
	s := &sm.steps[sm.flow.CurrentStep]
	s.client.Status = "done"
	s.resource.Status = "done"
	s.facilitator.Status = "done"
	sm.flow.CurrentStep++
	if sm.flow.CurrentStep < sm.flow.TotalSteps {
		sm.markRunning(&sm.steps[sm.flow.CurrentStep])
	}
}

func (sm *stepManager) markStepRunning() {
	if sm.flow.CurrentStep >= sm.flow.TotalSteps {
		return
	}
	s := &sm.steps[sm.flow.CurrentStep]
	if s.client.Action != inactiveAction {
		s.client.Status = "running"
	}
	if s.resource.Action != inactiveAction {
		s.resource.Status = "running"
	}
	if s.facilitator.Action != inactiveAction {
		s.facilitator.Status = "running"
	}
}

func (sm *stepManager) markStepError() {
	if sm.flow.CurrentStep >= sm.flow.TotalSteps {
		return
	}
	s := &sm.steps[sm.flow.CurrentStep]
	if s.client.Action != inactiveAction {
		s.client.Status = "error"
	}
	if s.resource.Action != inactiveAction {
		s.resource.Status = "error"
	}
	if s.facilitator.Action != inactiveAction {
		s.facilitator.Status = "error"
	}
}

func (sm *stepManager) view(width int) string {
	clientSteps := make([]StepInfo, 10)
	resourceSteps := make([]StepInfo, 10)
	facilitatorSteps := make([]StepInfo, 10)

	for i := range sm.steps {
		clientSteps[i] = sm.steps[i].client
		resourceSteps[i] = sm.steps[i].resource
		facilitatorSteps[i] = sm.steps[i].facilitator
	}

	return renderFlowPanels(
		sm.flow.CurrentStep, sm.flow.TotalSteps,
		clientSteps, resourceSteps, facilitatorSteps,
		width,
	)
}

func renderFlowPanels(
	step, totalSteps int,
	clientSteps, resourceSteps, facilitatorSteps []StepInfo,
	width int,
) string {
	progress := components.Progress{Total: totalSteps, Current: step}
	stepDesc := demo.StepDescription(step + 1)
	stepLabel := fmt.Sprintf("Step %d/%d: %s", step+1, totalSteps, stepDesc)

	header := lipgloss.JoinHorizontal(lipgloss.Center,
		progress.View(),
		"  ",
		lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).Render(stepLabel),
	)

	colWidth := max((width-panelGap)/3, minPanelWidth)

	clientPanel := renderActorPanel("Client", clientSteps, step, colWidth, tui.ColorSecondary)
	resourcePanel := renderActorPanel("Resource Server", resourceSteps, step, colWidth, tui.ColorPrimary)
	facilitatorPanel := renderActorPanel("Facilitator", facilitatorSteps, step, colWidth, tui.ColorAccent)

	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		clientPanel, " ", resourcePanel, " ", facilitatorPanel,
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, "", panels)
}

func renderActorPanel(title string, steps []StepInfo, currentStep, width int, color lipgloss.Color) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Width(width - 2).
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().Foreground(color).Bold(true)

	var b strings.Builder
	b.WriteString(titleStyle.Render(title) + "\n")

	for i, s := range steps {
		if i > currentStep+1 {
			break
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
