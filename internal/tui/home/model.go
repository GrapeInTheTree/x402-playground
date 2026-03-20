package home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

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

func New(width, height int) *Model {
	return &Model{
		menu:   components.NewMenu(menuItems),
		width:  width,
		height: height,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

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

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.menu.Width = width
}

func (m *Model) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		MarginLeft(2).
		Render("x402 Protocol Explorer")

	subtitle := lipgloss.NewStyle().
		Foreground(tui.ColorMuted).
		MarginLeft(2).
		MarginBottom(1).
		Render("Interactive learning tool for the x402 payment protocol")

	menu := m.menu.View()

	hints := components.StatusBar{Width: m.width}.View(
		"  ↑/↓ navigate  enter select  ? help  q quit",
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		subtitle,
		"",
		menu,
		hints,
	)
}
