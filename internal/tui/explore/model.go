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
	subPageERC8004
)

// Model is the explore page TUI model with sub-page navigation.
type Model struct {
	menu     components.Menu
	sub      subPage
	header   *HeaderModel
	typed    *TypedDataModel
	compare  *CompareModel
	onchain  *OnChainModel
	erc8004  *ERC8004Model
	width    int
	height   int
}

var menuItems = []components.MenuItem{
	{Title: "Decode PAYMENT-REQUIRED Header", Description: "Decode 402 response headers live", Icon: "📋"},
	{Title: "Inspect EIP-712 TypedData", Description: "Explore EIP-712 signature data structures", Icon: "🔬"},
	{Title: "Compare EIP-3009 vs Permit2", Description: "Side-by-side comparison of both methods", Icon: "⚖️"},
	{Title: "View On-Chain State", Description: "Balances, allowances, contract state", Icon: "🔗"},
	{Title: "ERC-8004 Agent Registries", Description: "Identity, Reputation, Validation for AI agents", Icon: "🤖"},
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
		erc8004: NewERC8004Model(width, height),
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
				case 4:
					m.sub = subPageERC8004
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
	case subPageERC8004:
		cmd = m.erc8004.Update(msg)
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
	m.erc8004.SetSize(width, height)
}

// View renders the current sub-page or the explore menu.
func (m *Model) View() string {
	availW := m.width - 4

	switch m.sub {
	case subPageMenu:
		return m.viewMenu(availW)
	case subPageHeader:
		return m.header.View()
	case subPageTypedData:
		return m.typed.View()
	case subPageCompare:
		return m.compare.View()
	case subPageOnChain:
		return m.onchain.View()
	case subPageERC8004:
		return m.erc8004.View()
	default:
		return ""
	}
}

func (m *Model) viewMenu(availW int) string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorBorder).
		Padding(0, 1)

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
		Render("Explore — Data Structure Inspector")
	subtitle := tui.MutedStyle.Render("Select a topic to explore x402 internals")

	boxW := min(availW-2, 64) // -2 for border
	content := lipgloss.JoinVertical(lipgloss.Left,
		title, subtitle, "", m.menu.View())

	return boxStyle.Width(boxW).Render(content)
}
