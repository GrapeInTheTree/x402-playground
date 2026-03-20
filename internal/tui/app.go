package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Page identifies which TUI page is active.
type Page int

const (
	PageHome Page = iota
	PageLearn
	PageExplore
	PagePractice
	PageDashboard
)

const (
	minTermWidth  = 60
	minTermHeight = 20
)

// SubModel is implemented by each page model.
type SubModel interface {
	Init() tea.Cmd
	Update(tea.Msg) (SubModel, tea.Cmd)
	View() string
	SetSize(width, height int)
}

// SubModelFactory creates a SubModel. Used to lazily initialize pages.
type SubModelFactory func(width, height int) SubModel

// RootModel is the top-level TUI model that routes between pages.
type RootModel struct {
	currentPage  Page
	pages        map[Page]SubModel
	factories    map[Page]SubModelFactory
	windowWidth  int
	windowHeight int
	showHelp     bool
}

func NewRootModel(factories map[Page]SubModelFactory) RootModel {
	return NewRootModelWithStart(factories, PageHome)
}

func NewRootModelWithStart(factories map[Page]SubModelFactory, startPage Page) RootModel {
	return RootModel{
		currentPage: startPage,
		pages:       make(map[Page]SubModel),
		factories:   factories,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		for page, sub := range m.pages {
			sub.SetSize(msg.Width, msg.Height-2)
			m.pages[page] = sub
		}
		return m, nil

	case NavigateMsg:
		m.currentPage = msg.Page
		if _, ok := m.pages[msg.Page]; !ok {
			if factory, ok := m.factories[msg.Page]; ok {
				sub := factory(m.windowWidth, m.windowHeight-2)
				m.pages[msg.Page] = sub
				return m, sub.Init()
			}
		}
		return m, nil

	case BackMsg:
		if m.currentPage != PageHome {
			m.currentPage = PageHome
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Toggle help overlay
		if msg.String() == "?" {
			m.showHelp = !m.showHelp
			return m, nil
		}
		// If help is showing, any other key closes it
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		// On home page, 'q' quits
		if m.currentPage == PageHome && msg.String() == "q" {
			return m, tea.Quit
		}
	}

	// Ensure current page is initialized
	if _, ok := m.pages[m.currentPage]; !ok {
		if factory, ok := m.factories[m.currentPage]; ok {
			sub := factory(m.windowWidth, m.windowHeight-2)
			m.pages[m.currentPage] = sub
			cmd := sub.Init()
			return m, cmd
		}
	}

	// Delegate to current page
	if sub, ok := m.pages[m.currentPage]; ok {
		newSub, cmd := sub.Update(msg)
		m.pages[m.currentPage] = newSub
		return m, cmd
	}

	return m, nil
}

func (m RootModel) View() string {
	if m.windowWidth == 0 {
		return "Loading..."
	}

	// Minimum terminal size check
	if m.windowWidth < minTermWidth || m.windowHeight < minTermHeight {
		msg := fmt.Sprintf(
			"Terminal too small: %dx%d\nMinimum required: %dx%d\n\nPlease resize your terminal.",
			m.windowWidth, m.windowHeight, minTermWidth, minTermHeight,
		)
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).
			Bold(true).
			Padding(2, 4).
			Render(msg)
	}

	// Render current page
	var content string
	if sub, ok := m.pages[m.currentPage]; ok {
		content = sub.View()
	} else {
		content = "Initializing..."
	}

	// Help overlay
	if m.showHelp {
		helpView := renderHelpOverlay(m.windowWidth, m.windowHeight)
		// Center the overlay
		helpView = lipgloss.Place(m.windowWidth, m.windowHeight, lipgloss.Center, lipgloss.Center, helpView)
		return helpView
	}

	return lipgloss.NewStyle().
		MaxWidth(m.windowWidth).
		MaxHeight(m.windowHeight).
		Render(content)
}

// renderHelpOverlay renders the keyboard shortcuts help panel inline
// to avoid a circular import with the components package.
func renderHelpOverlay(width, height int) string {
	title := lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Render("Keyboard Shortcuts")

	bindings := []struct{ key, desc string }{
		{"↑/k", "Move up"},
		{"↓/j", "Move down"},
		{"Enter", "Select / Confirm"},
		{"Esc", "Go back"},
		{"q", "Quit (from home)"},
		{"?", "Toggle this help"},
		{"", ""},
		{"n", "Next step (Practice)"},
		{"p", "Previous step (Practice)"},
		{"Tab", "Switch view (Explore)"},
		{"r", "Refresh (Dashboard)"},
	}

	var b strings.Builder
	b.WriteString(title + "\n\n")

	keyStyle := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Width(10)
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB"))

	for _, bind := range bindings {
		if bind.key == "" {
			b.WriteString("\n")
			continue
		}
		b.WriteString("  " + keyStyle.Render(bind.key) + descStyle.Render(bind.desc) + "\n")
	}

	b.WriteString("\n" + MutedStyle.Render("  Press ? to close"))

	boxWidth := 40
	if width > 0 && width < 50 {
		boxWidth = width - 10
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Width(boxWidth)

	return style.Render(b.String())
}
