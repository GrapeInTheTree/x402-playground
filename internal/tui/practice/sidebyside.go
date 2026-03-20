package practice

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/demo"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

// SideBySideModel runs EIP-3009 and Permit2 flows side by side.
type SideBySideModel struct {
	eip3009Step int
	permit2Step int
	totalSteps  int
	width       int
	height      int
	cfg         *config.ExplorerConfig
}

func NewSideBySideModel(width, height int, cfg *config.ExplorerConfig) *SideBySideModel {
	return &SideBySideModel{
		totalSteps: 10,
		width:      width,
		height:     height,
		cfg:        cfg,
	}
}

func (m *SideBySideModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "n":
			if m.eip3009Step < m.totalSteps-1 {
				m.eip3009Step++
			}
			if m.permit2Step < m.totalSteps-1 {
				m.permit2Step++
			}
		case "p":
			if m.eip3009Step > 0 {
				m.eip3009Step--
			}
			if m.permit2Step > 0 {
				m.permit2Step--
			}
		}
	}
	return nil
}

func (m *SideBySideModel) View() string {
	colWidth := (m.width - 6) / 2
	if colWidth < 30 {
		colWidth = 30
	}

	// Progress bars
	eipProgress := components.Progress{Total: m.totalSteps, Current: m.eip3009Step}
	p2Progress := components.Progress{Total: m.totalSteps, Current: m.permit2Step}

	leftTitle := lipgloss.NewStyle().
		Foreground(tui.ColorSecondary).
		Bold(true).
		Render("EIP-3009")

	rightTitle := lipgloss.NewStyle().
		Foreground(tui.ColorAccent).
		Bold(true).
		Render("Permit2")

	header := lipgloss.JoinHorizontal(lipgloss.Center,
		"  ", leftTitle, " ", eipProgress.View(),
		"      ",
		rightTitle, " ", p2Progress.View(),
	)

	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorSecondary).
		Width(colWidth - 2).
		Padding(1, 2)

	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorAccent).
		Width(colWidth - 2).
		Padding(1, 2)

	left := renderSideBySideSteps(m.eip3009Step, m.totalSteps, "eip3009", tui.ColorSecondary)
	right := renderSideBySideSteps(m.permit2Step, m.totalSteps, "permit2", tui.ColorAccent)

	panels := lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		"  ",
		rightStyle.Render(right),
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, "", panels)
}

func renderSideBySideSteps(currentStep, totalSteps int, method string, color lipgloss.Color) string {
	var b strings.Builder

	for i := 0; i <= currentStep && i < totalSteps; i++ {
		desc := demo.StepDescription(i + 1)
		var icon string
		var style lipgloss.Style

		if i < currentStep {
			icon = "✓"
			style = lipgloss.NewStyle().Foreground(tui.ColorSuccess)
		} else {
			icon = "►"
			style = lipgloss.NewStyle().Foreground(color).Bold(true)
		}

		b.WriteString(style.Render(fmt.Sprintf(" %s Step %d: %s", icon, i+1, desc)) + "\n")

		// Show method-specific detail for current step
		if i == currentStep {
			detail := stepDetail(i+1, method)
			if detail != "" {
				b.WriteString(tui.MutedStyle.Render("   "+detail) + "\n")
			}
		}
	}

	return b.String()
}

func stepDetail(step int, method string) string {
	if method == "permit2" {
		switch step {
		case 1:
			return "Permit2 approve 확인 필요"
		case 4:
			return "extra.assetTransferMethod: permit2"
		case 5:
			return "Domain: Permit2, Type: PermitWitnessTransferFrom"
		case 7:
			return "Permit2 서명 + allowance 검증"
		case 9:
			return "x402Permit2Proxy.settle() → Permit2"
		}
	} else {
		switch step {
		case 5:
			return "Domain: USDC, Type: TransferWithAuthorization"
		case 7:
			return "EIP-712 서명 + 잔액 + 시뮬레이션"
		case 9:
			return "USDC.transferWithAuthorization()"
		}
	}
	return ""
}
