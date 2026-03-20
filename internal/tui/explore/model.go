package explore

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

type subPage int

const (
	subPageMenu subPage = iota
	subPageHeader
	subPageTypedData
	subPageCompare
	subPageOnChain
)

// Model is the explore page TUI model with sub-page navigation.
type Model struct {
	menu     components.Menu
	sub      subPage
	header   *HeaderModel
	typed    *TypedDataModel
	compare  *CompareModel
	onchain  *OnChainModel
	width    int
	height   int
}

var menuItems = []components.MenuItem{
	{Title: "Decode PAYMENT-REQUIRED Header", Description: "402 응답 헤더 실시간 디코딩", Icon: "📋"},
	{Title: "Inspect EIP-712 TypedData", Description: "EIP-712 서명 데이터 구조 탐색", Icon: "🔬"},
	{Title: "Compare EIP-3009 vs Permit2", Description: "두 방식의 차이점 사이드바이사이드", Icon: "⚖️"},
	{Title: "View On-Chain State", Description: "잔액, allowance, 컨트랙트 상태", Icon: "🔗"},
}

// New creates a new explore page model with all sub-pages initialized.
func New(width, height int) *Model {
	return &Model{
		menu:    components.NewMenu(menuItems),
		sub:     subPageMenu,
		header:  NewHeaderModel(width, height),
		typed:   NewTypedDataModel(width, height),
		compare: NewCompareModel(width, height),
		onchain: NewOnChainModel(width, height),
		width:   width,
		height:  height,
	}
}

// Init implements the SubModel interface.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles key events for menu and sub-page navigation.
func (m *Model) Update(msg tea.Msg) (tui.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.sub == subPageMenu {
			switch msg.String() {
			case "esc", "q":
				return m, func() tea.Msg { return tui.BackMsg{} }
			case "up", "k":
				m.menu.Up()
			case "down", "j":
				m.menu.Down()
			case "enter":
				switch m.menu.Selected() {
				case 0:
					m.sub = subPageHeader
				case 1:
					m.sub = subPageTypedData
				case 2:
					m.sub = subPageCompare
				case 3:
					m.sub = subPageOnChain
				}
				return m, nil
			}
			return m, nil
		}

		// Sub-page handling
		if msg.String() == "esc" {
			m.sub = subPageMenu
			return m, nil
		}
	}

	// Delegate to active sub-page
	var cmd tea.Cmd
	switch m.sub {
	case subPageHeader:
		cmd = m.header.Update(msg)
	case subPageTypedData:
		cmd = m.typed.Update(msg)
	case subPageCompare:
		cmd = m.compare.Update(msg)
	case subPageOnChain:
		cmd = m.onchain.Update(msg)
	}

	return m, cmd
}

// SetSize updates dimensions for the model and all sub-pages.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.header.SetSize(width, height)
	m.typed.SetSize(width, height)
	m.compare.SetSize(width, height)
	m.onchain.SetSize(width, height)
}

// View renders the current sub-page or the explore menu.
func (m *Model) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		Render("Explore — Data Structure Inspector")

	var content string
	var hints string

	switch m.sub {
	case subPageMenu:
		content = m.menu.View()
		hints = "  ↑/↓ navigate  enter select  ? help  esc back"
	case subPageHeader:
		content = m.header.View()
		hints = "  ↑/↓ navigate fields  ? help  esc back to menu"
	case subPageTypedData:
		content = m.typed.View()
		hints = "  ↑/↓ navigate fields  ? help  esc back to menu"
	case subPageCompare:
		content = m.compare.View()
		hints = "  ? help  esc back to menu"
	case subPageOnChain:
		content = m.onchain.View()
		hints = "  ? help  esc back to menu"
	}

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		content,
	)

	return tui.LayoutPage(body, hints, m.width, m.height)
}
