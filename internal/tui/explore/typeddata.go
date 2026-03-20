package explore

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-demo/internal/tui"
	"github.com/GrapeInTheTree/x402-demo/internal/tui/components"
)

// TypedDataModel shows EIP-712 TypedData structure with field exploration.
type TypedDataModel struct {
	explorer  components.FieldExplorer
	showEIP   bool // false=EIP-3009, true=Permit2
	width     int
	height    int
}

func NewTypedDataModel(width, height int) TypedDataModel {
	return TypedDataModel{
		explorer: components.NewFieldExplorer(eip3009Fields()),
		width:    width,
		height:   height,
	}
}

func eip3009Fields() []components.Field {
	return []components.Field{
		{Name: "[Domain] name", Value: "USDC", Description: "EIP-712 도메인 이름. USDC 컨트랙트의 name() 반환값. Base Sepolia는 'USDC' (not 'USD Coin')."},
		{Name: "[Domain] version", Value: "2", Description: "EIP-712 도메인 버전. FiatTokenV2의 버전 문자열."},
		{Name: "[Domain] chainId", Value: "84532", Description: "Base Sepolia Chain ID. Replay protection에 사용."},
		{Name: "[Domain] contract", Value: "0x036CbD...", Description: "서명을 검증하는 컨트랙트 주소 (USDC). Domain separation의 핵심."},
		{Name: "[Msg] from", Value: "0xClient...", Description: "서명자이자 USDC 보유자. 서명 복원으로 이 주소 일치를 검증."},
		{Name: "[Msg] to", Value: "0xPayTo...", Description: "USDC 수신자 주소. PAYMENT-REQUIRED의 payTo와 일치해야 함."},
		{Name: "[Msg] value", Value: "100000", Description: "전송 금액 (0.1 USDC). maxAmountRequired 이상이어야 검증 통과."},
		{Name: "[Msg] validAfter", Value: "0", Description: "서명 유효 시작 시간 (Unix). 0 = 즉시 유효."},
		{Name: "[Msg] validBefore", Value: "1718000000", Description: "서명 만료 시간 (Unix). 이 시간 이후에는 트랜잭션 실행 불가."},
		{Name: "[Msg] nonce", Value: "0xabcd...1234", Description: "랜덤 32바이트. 이중 사용 방지 — 한번 사용된 nonce는 재사용 불가."},
	}
}

func permit2Fields() []components.Field {
	return []components.Field{
		{Name: "[Domain] name", Value: "Permit2", Description: "Permit2 컨트랙트의 EIP-712 도메인 이름."},
		{Name: "[Domain] chainId", Value: "84532", Description: "Base Sepolia Chain ID."},
		{Name: "[Domain] contract", Value: "0x000...22D4", Description: "Permit2 컨트랙트 주소 (CREATE2, 모든 체인 동일)."},
		{Name: "[Msg] permitted.token", Value: "0x036CbD...", Description: "전송할 토큰 주소 (USDC). EIP-3009와 달리 모든 ERC-20 가능."},
		{Name: "[Msg] permitted.amount", Value: "100000", Description: "전송 허용 금액 (0.1 USDC)."},
		{Name: "[Msg] spender", Value: "0x4020...0001", Description: "x402Permit2Proxy 주소. Permit2가 이 주소에게 토큰 이동을 허용."},
		{Name: "[Msg] nonce", Value: "12345", Description: "Permit2 nonce. EIP-3009의 랜덤 nonce와 달리, 순차적으로 증가."},
		{Name: "[Msg] deadline", Value: "1718000000", Description: "서명 만료 시간 (Unix)."},
		{Name: "[Witness] to", Value: "0xPayTo...", Description: "최종 토큰 수신자. x402Permit2Proxy가 이 주소로 토큰 전송."},
		{Name: "[Witness] validAfter", Value: "0", Description: "유효 시작 시간. EIP-3009와 동일한 역할."},
	}
}

func (m *TypedDataModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			m.explorer.Up()
		case "down", "j":
			m.explorer.Down()
		case "tab":
			m.showEIP = !m.showEIP
			if m.showEIP {
				m.explorer = components.NewFieldExplorer(permit2Fields())
			} else {
				m.explorer = components.NewFieldExplorer(eip3009Fields())
			}
		}
	}
	return nil
}

func (m *TypedDataModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.explorer.Width = width
}

func (m *TypedDataModel) View() string {
	mode := "EIP-3009 (USDC Direct)"
	if m.showEIP {
		mode = "Permit2 (Universal ERC-20)"
	}

	title := lipgloss.NewStyle().
		Foreground(tui.ColorSecondary).
		Bold(true).
		MarginLeft(4).
		Render("EIP-712 TypedData — " + mode)

	hint := tui.MutedStyle.
		MarginLeft(4).
		Render("Tab으로 EIP-3009/Permit2 전환  ↑/↓ 필드 선택")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		hint,
		"",
		m.explorer.View(),
	)
}
