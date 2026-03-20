package home

import (
	"fmt"
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
	{Title: "Learn", Description: "Learn x402 protocol with coding quizzes", Icon: "\u25c8"},
	{Title: "Explore", Description: "Inspect protocol data structures live", Icon: "\u25ce"},
	{Title: "Practice", Description: "Execute payment flows (EIP-3009 / Permit2)", Icon: "\u25b6"},
	{Title: "Dashboard", Description: "Wallet balances & transaction status", Icon: "\u25eb"},
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
	contentWidth := min(m.width-8, 60)

	// ASCII art banner
	bannerLine := strings.Repeat("\u2550", contentWidth-2)
	banner := fmt.Sprintf("  %s\n   x402 Protocol Explorer\n  %s", bannerLine, bannerLine)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		Render(banner)

	subtitle := lipgloss.NewStyle().
		Foreground(tui.ColorMuted).
		Render("  Interactive learning tool for the x402 payment protocol")

	body := lipgloss.JoinVertical(lipgloss.Left,
		"",
		title,
		subtitle,
		"",
		"",
		m.menu.View(),
	)

	return tui.LayoutPageCentered(body, "\u2191/\u2193 navigate  enter select  ? help  q quit", m.width, m.height)
}
