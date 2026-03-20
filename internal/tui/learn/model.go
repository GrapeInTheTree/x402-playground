package learn

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

type view int

const (
	viewTopicList view = iota
	viewContent
)

type Model struct {
	topics   []Topic
	menu     components.Menu
	viewport viewport.Model
	current  view
	width    int
	height   int
}

func New(width, height int) *Model {
	topics := AllTopics()
	items := make([]components.MenuItem, len(topics))
	for i, t := range topics {
		items[i] = components.MenuItem{
			Title:       t.Title,
			Description: t.Description,
			Icon:        t.Icon,
		}
	}

	vp := viewport.New(width-4, height-6)
	vp.SetContent("")

	return &Model{
		topics:   topics,
		menu:     components.NewMenu(items),
		viewport: vp,
		current:  viewTopicList,
		width:    width,
		height:   height,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tui.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.current == viewContent {
				m.current = viewTopicList
				return m, nil
			}
			return m, func() tea.Msg { return tui.BackMsg{} }
		case "q":
			if m.current == viewTopicList {
				return m, func() tea.Msg { return tui.BackMsg{} }
			}
		case "up", "k":
			if m.current == viewTopicList {
				m.menu.Up()
			}
		case "down", "j":
			if m.current == viewTopicList {
				m.menu.Down()
			}
		case "enter":
			if m.current == viewTopicList {
				idx := m.menu.Selected()
				if idx >= 0 && idx < len(m.topics) {
					content := components.RenderMarkdown(
						m.topics[idx].Content(),
						m.width-8,
					)
					m.viewport.SetContent(content)
					m.viewport.GotoTop()
					m.current = viewContent
				}
				return m, nil
			}
		}
	}

	if m.current == viewContent {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width - 4
	m.viewport.Height = height - 6
	m.menu.Width = width
}

func (m *Model) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		MarginLeft(2).
		Render("Learn — x402 Protocol")

	var content string
	var hints string

	if m.current == viewTopicList {
		content = m.menu.View()
		hints = "  ↑/↓ navigate  enter select  ? help  esc back"
	} else {
		scrollPct := fmt.Sprintf("%d%%", int(m.viewport.ScrollPercent()*100))
		title := lipgloss.NewStyle().
			Foreground(tui.ColorSecondary).
			MarginLeft(4).
			Render(m.topics[m.menu.Selected()].Title + " " +
				tui.MutedStyle.Render("["+scrollPct+"]"))
		content = title + "\n" + m.viewport.View()
		hints = "  ↑/↓ scroll  ? help  esc back to topics"
	}

	statusBar := components.StatusBar{Width: m.width}.View(hints)

	return lipgloss.JoinVertical(lipgloss.Left,
		"",
		header,
		"",
		content,
		statusBar,
	)
}
