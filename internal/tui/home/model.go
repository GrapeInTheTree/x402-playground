package home

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

// Model is the home page TUI model with a navigation menu.
type Model struct {
	menu   components.Menu
	width  int
	height int
}

var menuItems = []components.MenuItem{
	{Title: "Learn", Description: "x402 프로토콜 개념 학습", Icon: "📖"},
	{Title: "Explore", Description: "실시간 데이터 구조 탐색", Icon: "🔍"},
	{Title: "Practice", Description: "결제 흐름 실행 (EIP-3009 / Permit2)", Icon: "⚡"},
	{Title: "Dashboard", Description: "지갑 잔액 & 트랜잭션 현황", Icon: "📊"},
}

var pageMap = []tui.Page{
	tui.PageLearn,
	tui.PageExplore,
	tui.PagePractice,
	tui.PageDashboard,
}

// New creates a new home page model with the given dimensions.
func New(width, height int) *Model {
	return &Model{
		menu:   components.NewMenu(menuItems),
		width:  width,
		height: height,
	}
}

// Init implements the SubModel interface.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles key events for menu navigation.
func (m *Model) Update(msg tea.Msg) (tui.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.menu.Up()
		case "down", "j":
			m.menu.Down()
		case "enter":
			idx := m.menu.Selected()
			if idx >= 0 && idx < len(pageMap) {
				return m, func() tea.Msg {
					return tui.NavigateMsg{Page: pageMap[idx]}
				}
			}
		}
	}
	return m, nil
}

// SetSize updates the model dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.menu.Width = width
}

// View renders the home page with title and menu.
func (m *Model) View() string {
	// Header section
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		Render("x402 Protocol Explorer")

	subtitle := lipgloss.NewStyle().
		Foreground(tui.ColorMuted).
		Render("Interactive learning tool for the x402 payment protocol")

	divider := lipgloss.NewStyle().
		Foreground(tui.ColorBorder).
		Width(min(m.width-8, 60)).
		Render(strings.Repeat("─", min(m.width-8, 60)))

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		title,
		subtitle,
		"",
		divider,
		"",
		m.menu.View(),
	)

	return tui.LayoutPage(body, "↑/↓ navigate  enter select  ? help  q quit", m.width, m.height)
}
