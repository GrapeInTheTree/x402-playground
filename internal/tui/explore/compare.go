package explore

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// CompareModel shows EIP-3009 vs Permit2 side-by-side.
type CompareModel struct {
	width  int
	height int
}

func NewCompareModel(width, height int) CompareModel {
	return CompareModel{width: width, height: height}
}

func (m *CompareModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (m *CompareModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *CompareModel) View() string {
	colWidth := (m.width - 10) / 2
	if colWidth < 30 {
		colWidth = 30
	}

	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorSecondary).
		Width(colWidth).
		Padding(1, 2)

	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorAccent).
		Width(colWidth).
		Padding(1, 2)

	left := buildEIP3009Panel()
	right := buildPermit2Panel()

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		"  ",
		rightStyle.Render(right),
	)
}

func buildEIP3009Panel() string {
	var b strings.Builder
	title := lipgloss.NewStyle().
		Foreground(tui.ColorSecondary).
		Bold(true).
		Render("EIP-3009")

	b.WriteString(title + "\n")
	b.WriteString(tui.MutedStyle.Render("transferWithAuthorization") + "\n\n")

	rows := []struct{ key, val string }{
		{"Token Support", "EIP-3009 구현 토큰만 (USDC, EURC)"},
		{"Prerequisites", "없음"},
		{"Domain", "USDC 컨트랙트"},
		{"Primary Type", "TransferWithAuthorization"},
		{"Contracts", "토큰 컨트랙트 직접 호출"},
		{"Nonce", "랜덤 32바이트 (1회용)"},
		{"On-chain Call", "USDC.transferWithAuthorization(...)"},
		{"Gas", "Facilitator 대납"},
	}

	for _, r := range rows {
		key := lipgloss.NewStyle().
			Foreground(tui.ColorSecondary).
			Bold(true).
			Render(r.key + ":")
		b.WriteString(key + "\n")
		b.WriteString("  " + r.val + "\n\n")
	}

	return b.String()
}

func buildPermit2Panel() string {
	var b strings.Builder
	title := lipgloss.NewStyle().
		Foreground(tui.ColorAccent).
		Bold(true).
		Render("Permit2")

	b.WriteString(title + "\n")
	b.WriteString(tui.MutedStyle.Render("permitWitnessTransferFrom") + "\n\n")

	rows := []struct{ key, val string }{
		{"Token Support", "모든 ERC-20 토큰"},
		{"Prerequisites", "approve(Permit2, amount) 1회"},
		{"Domain", "Permit2 컨트랙트"},
		{"Primary Type", "PermitWitnessTransferFrom"},
		{"Contracts", "Permit2 + x402Permit2Proxy"},
		{"Nonce", "순차적 Permit2 nonce"},
		{"On-chain Call", "Proxy.settle(...) → Permit2"},
		{"Gas", "Facilitator 대납"},
	}

	for _, r := range rows {
		key := lipgloss.NewStyle().
			Foreground(tui.ColorAccent).
			Bold(true).
			Render(r.key + ":")
		b.WriteString(key + "\n")
		b.WriteString("  " + r.val + "\n\n")
	}

	return b.String()
}
