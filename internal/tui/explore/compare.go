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

// NewCompareModel creates a new side-by-side comparison model.
func NewCompareModel(width, height int) *CompareModel {
	return &CompareModel{width: width, height: height}
}

// Update is a no-op since the compare view is static.
func (m *CompareModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// SetSize updates the model dimensions.
func (m *CompareModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

const minCompareColWidth = 30

// View renders the EIP-3009 vs Permit2 side-by-side comparison.
func (m *CompareModel) View() string {
	colWidth := max((m.width-10)/2, minCompareColWidth)

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

	left := buildComparePanel("EIP-3009", "transferWithAuthorization", tui.ColorSecondary, eip3009Rows)
	right := buildComparePanel("Permit2", "permitWitnessTransferFrom", tui.ColorAccent, permit2Rows)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		"  ",
		rightStyle.Render(right),
	)
}

type compareRow struct{ Key, Val string }

var eip3009Rows = []compareRow{
	{"Token Support", "EIP-3009 구현 토큰만 (USDC, EURC)"},
	{"Prerequisites", "없음"},
	{"Domain", "USDC 컨트랙트"},
	{"Primary Type", "TransferWithAuthorization"},
	{"Contracts", "토큰 컨트랙트 직접 호출"},
	{"Nonce", "랜덤 32바이트 (1회용)"},
	{"On-chain Call", "USDC.transferWithAuthorization(...)"},
	{"Gas", "Facilitator 대납"},
}

var permit2Rows = []compareRow{
	{"Token Support", "모든 ERC-20 토큰"},
	{"Prerequisites", "approve(Permit2, amount) 1회"},
	{"Domain", "Permit2 컨트랙트"},
	{"Primary Type", "PermitWitnessTransferFrom"},
	{"Contracts", "Permit2 + x402Permit2Proxy"},
	{"Nonce", "순차적 Permit2 nonce"},
	{"On-chain Call", "Proxy.settle(...) → Permit2"},
	{"Gas", "Facilitator 대납"},
}

func buildComparePanel(title, subtitle string, color lipgloss.Color, rows []compareRow) string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Foreground(color).Bold(true).Render(title) + "\n")
	b.WriteString(tui.MutedStyle.Render(subtitle) + "\n\n")

	keyStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
	for _, r := range rows {
		b.WriteString(keyStyle.Render(r.Key+":") + "\n")
		b.WriteString("  " + r.Val + "\n\n")
	}

	return b.String()
}
