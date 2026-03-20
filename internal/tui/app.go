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

// pageLabel returns the display name for each page.
func pageLabel(p Page) string {
	switch p {
	case PageHome:
		return "Home"
	case PageLearn:
		return "Learn"
	case PageExplore:
		return "Explore"
	case PagePractice:
		return "Practice"
	case PageDashboard:
		return "Dashboard"
	default:
		return ""
	}
}

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

// NewRootModel creates a new root model starting at the home page.
func NewRootModel(factories map[Page]SubModelFactory) RootModel {
	return NewRootModelWithStart(factories, PageHome)
}

// NewRootModelWithStart creates a new root model starting at the specified page.
func NewRootModelWithStart(factories map[Page]SubModelFactory, startPage Page) RootModel {
	return RootModel{
		currentPage: startPage,
		pages:       make(map[Page]SubModel),
		factories:   factories,
	}
}

// Init implements the tea.Model interface.
func (m RootModel) Init() tea.Cmd {
	return nil
}

// Update handles window resize, navigation, help toggle, and delegates to the active page.
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// Pages get the inner content size (frame takes some space)
		iw, ih := m.innerSize()
		for page, sub := range m.pages {
			sub.SetSize(iw, ih)
			m.pages[page] = sub
		}
		if _, ok := m.pages[m.currentPage]; !ok {
			if factory, ok := m.factories[m.currentPage]; ok {
				sub := factory(iw, ih)
				m.pages[m.currentPage] = sub
				return m, sub.Init()
			}
		}
		return m, nil

	case NavigateMsg:
		m.currentPage = msg.Page
		if _, ok := m.pages[msg.Page]; !ok {
			if factory, ok := m.factories[msg.Page]; ok {
				iw, ih := m.innerSize()
				sub := factory(iw, ih)
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
		if msg.String() == "?" {
			m.showHelp = !m.showHelp
			return m, nil
		}
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		if m.currentPage == PageHome && msg.String() == "q" {
			return m, tea.Quit
		}
	}

	if _, ok := m.pages[m.currentPage]; !ok {
		if factory, ok := m.factories[m.currentPage]; ok {
			iw, ih := m.innerSize()
			sub := factory(iw, ih)
			m.pages[m.currentPage] = sub
			return m, sub.Init()
		}
	}

	if sub, ok := m.pages[m.currentPage]; ok {
		newSub, cmd := sub.Update(msg)
		m.pages[m.currentPage] = newSub
		return m, cmd
	}

	return m, nil
}

// innerSize returns the content area size inside the app frame.
// Frame: 2 border + 2 padding horizontal, header + status + border vertical.
func (m RootModel) innerSize() (int, int) {
	iw := max(m.windowWidth-6, 20)  // 2 border + 4 padding
	ih := max(m.windowHeight-6, 10) // header(1) + border(2) + status(1) + padding(2)
	return iw, ih
}

// View renders the app frame with header bar, page content, and status bar.
func (m RootModel) View() string {
	w := m.windowWidth
	h := m.windowHeight

	if w == 0 {
		return "Loading..."
	}

	if w < minTermWidth || h < minTermHeight {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")).Bold(true).Padding(2, 4).
			Render(fmt.Sprintf("Terminal too small: %dx%d\nMinimum: %dx%d\n\nPlease resize.", w, h, minTermWidth, minTermHeight))
	}

	// Help overlay
	if m.showHelp {
		helpView := renderHelpOverlay(w)
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, helpView)
	}

	// Get page content
	var pageContent string
	if sub, ok := m.pages[m.currentPage]; ok {
		pageContent = sub.View()
	} else {
		pageContent = "Initializing..."
	}

	// === Build the consistent app frame ===

	innerW := max(w-4, 20) // border takes 2 chars each side

	// Header bar: app name (left) + current page tab (right)
	appName := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Render(" x402 Playground")

	pageTab := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true).
		Render(pageLabel(m.currentPage) + " ")

	spacerW := max(innerW-lipgloss.Width(appName)-lipgloss.Width(pageTab), 0)
	headerBar := appName + strings.Repeat(" ", spacerW) + pageTab

	headerStyle := lipgloss.NewStyle().
		Background(ColorSubtle).
		Width(innerW).
		Padding(0, 1)

	header := headerStyle.Render(headerBar)

	// Status bar (inside the frame, at bottom)
	hints := m.statusHints()
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Background(ColorSubtle).
		Width(innerW).
		Padding(0, 1)
	status := statusStyle.Render(" " + hints)

	// Content area: fill remaining height
	headerH := lipgloss.Height(header)
	statusH := lipgloss.Height(status)
	borderH := 2 // top + bottom border
	contentH := max(h-headerH-statusH-borderH, 1)

	content := lipgloss.NewStyle().
		Width(innerW).
		Height(contentH).
		Padding(0, 1).
		Render(pageContent)

	// Assemble inside frame
	inside := lipgloss.JoinVertical(lipgloss.Top, header, content, status)

	// Wrap in border
	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Width(w - 2).
		Render(inside)

	return frame
}

// statusHints returns context-sensitive keyboard hints for the current page.
func (m RootModel) statusHints() string {
	switch m.currentPage {
	case PageHome:
		return "↑/↓ navigate  enter select  ? help  q quit"
	case PageLearn:
		return "↑/↓ navigate  enter select/edit  ? help  esc back"
	case PageExplore:
		return "↑/↓ navigate  tab switch  ? help  esc back"
	case PagePractice:
		return "n next step  p prev  ? help  esc back"
	case PageDashboard:
		return "r refresh  ? help  esc back"
	default:
		return "? help  esc back"
	}
}

// renderHelpOverlay renders the keyboard shortcuts help panel.
func renderHelpOverlay(width int) string {
	title := lipgloss.NewStyle().
		Foreground(ColorPrimary).Bold(true).
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
		{"e", "Open editor (Learn)"},
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

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Width(min(40, width-10)).
		Render(b.String())
}
