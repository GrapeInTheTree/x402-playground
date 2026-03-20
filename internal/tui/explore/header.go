package explore

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-demo/internal/tui"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/components"
)

// HeaderModel shows the PAYMENT-REQUIRED header structure with field exploration.
type HeaderModel struct {
	explorer components.FieldExplorer
	width    int
	height   int
}

func NewHeaderModel(width, height int) HeaderModel {
	fields := []components.Field{
		{Name: "scheme", Value: "exact", Description: "결제 방식. 'exact'는 정확한 금액 결제를 의미. 현재 x402에서 유일하게 지원되는 scheme."},
		{Name: "network", Value: "eip155:84532", Description: "CAIP-2 네트워크 식별자. 'eip155'는 EVM 체인, '84532'는 Base Sepolia의 Chain ID."},
		{Name: "maxAmountRequired", Value: "100000", Description: "최대 결제 금액 (smallest unit). USDC는 6 decimals → 100000 = 0.1 USDC = $0.10"},
		{Name: "resource", Value: "https://.../ weather", Description: "결제 대상 리소스 URL. Client가 접근하려는 API 엔드포인트."},
		{Name: "payTo", Value: "0x1234...abcd", Description: "USDC를 수신할 주소. Private key가 필요하지 않은 수신 전용 주소."},
		{Name: "asset", Value: "0x036CbD...3dCF7e", Description: "결제에 사용할 ERC-20 토큰 컨트랙트 주소 (Base Sepolia USDC)."},
		{Name: "extra.name", Value: "USDC", Description: "EIP-712 도메인의 name 필드. 토큰 컨트랙트의 name() 반환값과 정확히 일치해야 한다."},
		{Name: "extra.version", Value: "2", Description: "EIP-712 도메인의 version 필드. USDC v2 (FiatTokenV2) 컨트랙트."},
		{Name: "extra.assetTransferMethod", Value: "(optional)", Description: "'permit2'이면 Permit2 방식. 미지정 시 EIP-3009 기본. Client SDK가 이 값으로 서명 방식 결정."},
	}

	return HeaderModel{
		explorer: components.NewFieldExplorer(fields),
		width:    width,
		height:   height,
	}
}

func (m *HeaderModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			m.explorer.Up()
		case "down", "j":
			m.explorer.Down()
		}
	}
	return nil
}

func (m *HeaderModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.explorer.Width = width
}

func (m *HeaderModel) View() string {
	title := lipgloss.NewStyle().
		Foreground(tui.ColorSecondary).
		Bold(true).
		MarginLeft(4).
		Render("PAYMENT-REQUIRED Header Fields")

	subtitle := tui.MutedStyle.
		MarginLeft(4).
		Render("↑/↓로 필드 선택 — 각 필드의 역할과 의미를 확인하세요")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		subtitle,
		"",
		m.explorer.View(),
	)
}
