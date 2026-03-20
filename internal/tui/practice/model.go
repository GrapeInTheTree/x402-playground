package practice

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/config"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

type subPage int

const (
	subPageMenu subPage = iota
	subPageEIP3009
	subPagePermit2
	subPageSideBySide
)

// Model is the practice page TUI model with payment flow sub-pages.
type Model struct {
	menu        components.Menu
	sub         subPage
	eip3009flow *EIP3009FlowModel
	permit2flow *Permit2FlowModel
	sidebyside  *SideBySideModel
	cfg         *config.ExplorerConfig
	width       int
	height      int
}

var menuItems = []components.MenuItem{
	{Title: "EIP-3009 Full Flow", Description: "10-step 결제 흐름 (USDC transferWithAuthorization)", Icon: "✍️"},
	{Title: "Permit2 Full Flow", Description: "10-step 결제 흐름 (Permit2 경유)", Icon: "🔑"},
	{Title: "Side-by-Side Comparison", Description: "EIP-3009과 Permit2 동시 실행 비교", Icon: "⚖️"},
}

// New creates a new practice page model with the given dimensions and configuration.
func New(width, height int, cfg *config.ExplorerConfig) *Model {
	return &Model{
		menu:   components.NewMenu(menuItems),
		sub:    subPageMenu,
		cfg:    cfg,
		width:  width,
		height: height,
	}
}

// NewWithFlow creates a Practice model that auto-starts a specific flow.
func NewWithFlow(width, height int, cfg *config.ExplorerConfig, flow string) *Model {
	m := New(width, height, cfg)
	switch flow {
	case "eip3009":
		m.eip3009flow = NewEIP3009FlowModel(width, height, cfg)
		m.sub = subPageEIP3009
	case "permit2":
		m.permit2flow = NewPermit2FlowModel(width, height, cfg)
		m.sub = subPagePermit2
	case "sidebyside":
		m.sidebyside = NewSideBySideModel(width, height, cfg)
		m.sub = subPageSideBySide
	}
	return m
}

// Init starts the spinner for any pre-selected sub-flow.
func (m *Model) Init() tea.Cmd {
	// If a sub-flow was pre-selected (via NewWithFlow), start its spinner
	switch m.sub {
	case subPageEIP3009:
		if m.eip3009flow != nil {
			return m.eip3009flow.Init()
		}
	case subPagePermit2:
		if m.permit2flow != nil {
			return m.permit2flow.Init()
		}
	}
	return nil
}

// Update handles key events for menu navigation and delegates to the active sub-flow.
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
				var initCmd tea.Cmd
				switch m.menu.Selected() {
				case 0:
					m.eip3009flow = NewEIP3009FlowModel(m.width, m.height, m.cfg)
					m.sub = subPageEIP3009
					initCmd = m.eip3009flow.Init()
				case 1:
					m.permit2flow = NewPermit2FlowModel(m.width, m.height, m.cfg)
					m.sub = subPagePermit2
					initCmd = m.permit2flow.Init()
				case 2:
					m.sidebyside = NewSideBySideModel(m.width, m.height, m.cfg)
					m.sub = subPageSideBySide
				}
				return m, initCmd
			}
			return m, nil
		}

		if msg.String() == "esc" {
			m.sub = subPageMenu
			return m, nil
		}
	}

	// Delegate to active sub-page
	var cmd tea.Cmd
	switch m.sub {
	case subPageEIP3009:
		if m.eip3009flow != nil {
			cmd = m.eip3009flow.Update(msg)
		}
	case subPagePermit2:
		if m.permit2flow != nil {
			cmd = m.permit2flow.Update(msg)
		}
	case subPageSideBySide:
		if m.sidebyside != nil {
			cmd = m.sidebyside.Update(msg)
		}
	}

	return m, cmd
}

// SetSize updates the model dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the current sub-page or the practice menu.
func (m *Model) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		Render("Practice — Payment Flow Execution")

	var content string
	var hints string

	switch m.sub {
	case subPageMenu:
		content = m.menu.View()
		hints = "  ↑/↓ navigate  enter select  ? help  esc back"
	case subPageEIP3009:
		if m.eip3009flow != nil {
			content = m.eip3009flow.View()
		}
		hints = "  n next step  p prev step  ? help  esc back to menu"
	case subPagePermit2:
		if m.permit2flow != nil {
			content = m.permit2flow.View()
		}
		hints = "  n next step  p prev step  ? help  esc back to menu"
	case subPageSideBySide:
		if m.sidebyside != nil {
			content = m.sidebyside.View()
		}
		hints = "  n next step  ? help  esc back to menu"
	}

	statusBar := components.StatusBar{Width: m.width}.View(hints)

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		content,
	)

	centered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, body)

	return lipgloss.JoinVertical(lipgloss.Left,
		centered,
		statusBar,
	)
}
